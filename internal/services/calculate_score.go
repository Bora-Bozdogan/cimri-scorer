package services

import (
	"scorer/internal/models"
)

type merchantProductRepoInterface interface {
	Read(merchantID int, productID int) *models.MerchantProduct
}

type productRepoInteface interface {
	ReadID(ID int) *models.Product
	ReadName(name string) *models.Product
}

type merchantRepoInterface interface {
	ReadName(name string) *models.Merchant
}

type queServiceClient interface {
	EnqueueRequest(m models.Request, score int) error
}

type metricsInterface interface {
	IncrementRequestCount()
	IncrementValidRequestCount()
}

type ServicesFuncs struct {
	productRepo         productRepoInteface
	merchantRepo        merchantRepoInterface
	merchantProductRepo merchantProductRepoInterface
	queServiceClient    queServiceClient
	metrics metricsInterface
}

func NewServicesFuncs(productRepo productRepoInteface, merchantRepo merchantRepoInterface, merchantProductRepo merchantProductRepoInterface, queServiceClient queServiceClient, metrics metricsInterface) ServicesFuncs {
	return ServicesFuncs{productRepo: productRepo, merchantRepo: merchantRepo, merchantProductRepo: merchantProductRepo, queServiceClient: queServiceClient, metrics: metrics}
}

func (s ServicesFuncs) CalculateScore(req *models.Request) int {
	//store tier
	storeTier := 1
	tiers := map[string]int{
		"Amazon":      10,
		"Ebay":        9,
		"Trendyol":    8,
		"Hepsiburada": 7,
		"Bolgesel":    5,
		"Yerel":       3,
		"Test":        1,
	}
	storeTier = tiers[*req.StoreName]

	//update type
	updateType := 0

	/*
		to compute updateType, we need to access past product desc, popularity, urgency -> relevant product db
		we also need past price for merchant, past stock -> merchantProduct

		first, get product struct matching name of request
		next, get merchant struct matching name of store
		finally, get productMerchant struct matching productID and merchantID
	*/
	product := s.productRepo.ReadName(*req.ProductName)
	merchant := s.merchantRepo.ReadName(*req.StoreName)
	var merchantProduct *models.MerchantProduct
	if product == nil || merchant == nil {
		merchantProduct = nil
	} else {
		merchantProduct = s.merchantProductRepo.Read(merchant.ID, product.ID)
	}
	if merchantProduct == nil {
		//update type is a 7, new item
		updateType = 7
	} else {
		priceChanged := float64(merchantProduct.MerchantPrice-*req.Price) / float64(*req.Price)
		if priceChanged > 0.2 {
			//check +%20 price drop
			updateType = 10
		} else if (*req.Stock == 0 && merchantProduct.MerchantStock != 0) ||
			(*req.Stock != 0 && merchantProduct.MerchantStock == 0) ||
			(*req.Stock < 5 && merchantProduct.MerchantStock > 5) {
			//check if critical stock change
			updateType = 9
		} else if priceChanged >= 0.05 {
			//check if important price change (%5-%20)
			updateType = 8
		} else if priceChanged > 0 {
			//unimportant price change (%5-%20)
			updateType = 5
		} else {
			//product description, or image change
			//get product for this
			/* can't throw error, because if merchant product exists,
			product referenced by it must also exist, otherwise psql would throw error */
			if *req.ProductDescription != product.ProductDescription ||
				*req.ProductImage != product.ProductImage {
				updateType = 3
			}
		}
	}
	productPopularity := 0
	urgencyFactor := 0

	if product == nil {
		//get popularity score
		productPopularity = *req.PopularityScore

		//get urgency score
		urgencyFactor = *req.UrgencyScore
	} else {
		//get popularity score
		productPopularity = product.Popularity_score

		//get urgency score
		urgencyFactor = product.Urgency_score
	}

	//calculate time decay
	timeDecay := 0 //2 points per minute, check que for last added same item, write post que implementation

	//final score
	score := (storeTier * 30) +
		(updateType * 25) +
		(productPopularity * 20) +
		(urgencyFactor * 15) +
		(timeDecay * 10)

	return score
}

func (s ServicesFuncs) EnqueueRequest(m models.Request, score int) error {
	return s.queServiceClient.EnqueueRequest(m, score)
}

func (s ServicesFuncs) IncrementRequestCount() {
	s.metrics.IncrementRequestCount()
}

func (s ServicesFuncs) IncrementValidRequestCount() {
	s.metrics.IncrementValidRequestCount()
}