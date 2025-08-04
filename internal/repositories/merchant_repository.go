package repositories

import (
	"scorer/internal/models"
	"gorm.io/gorm"
)

// struct
type MerchantRepository struct {
	db *gorm.DB
}

// constructor
func NewMerchantRepository(db *gorm.DB) MerchantRepository {
	repo := new(MerchantRepository)
	repo.db = db
	return *repo
}

// read command
func (r MerchantRepository) ReadName(name string) *models.Merchant {
	res := &models.Merchant{}
	err := r.db.First(res, "name = ?", name)
	if err.Error != nil {
		return nil
	}
	return res
}
