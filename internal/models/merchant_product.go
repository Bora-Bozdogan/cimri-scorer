package models

type MerchantProduct struct {
	ID int
	ProductID int
	MerchantID int 
	MerchantPrice int `gorm:"column:price"`
	MerchantStock int `gorm:"column:stock"`
}