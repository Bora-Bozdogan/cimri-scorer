package repositories

import (
	"scorer/internal/models"
	"gorm.io/gorm"
)

// struct
type MerchantProductRepository struct {
	db *gorm.DB
}

// constructor
func NewMerchantProductRepository(db *gorm.DB) MerchantProductRepository {
	repo := new(MerchantProductRepository)
	repo.db = db
	return *repo
}

// read command
func (r MerchantProductRepository) Read(merchantID int, productID int) *models.MerchantProduct {
	res := &models.MerchantProduct{}
	err := r.db.First(res, "merchant_id = ? AND product_id = ?", merchantID, productID)
	if err.Error != nil {
		return nil
	}
	return res
}
