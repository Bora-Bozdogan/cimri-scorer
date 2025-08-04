package repositories

import (
	"scorer/internal/models"
	"gorm.io/gorm"
)

// struct
type ProductRepository struct {
	db *gorm.DB
}

// constructor
func NewProductRepository(db *gorm.DB) ProductRepository {
	repo := new(ProductRepository)
	repo.db = db
	return *repo
}

// read command
func (r ProductRepository) ReadID(ID int) *models.Product {
	res := &models.Product{}
	err := r.db.First(res, "id = ?", ID)
	if err.Error != nil {
		return nil
	}
	return res
}

func (r ProductRepository) ReadName(name string) *models.Product {
	res := &models.Product{}
	err := r.db.First(res, "name = ?", name)
	if err.Error != nil {
		return nil
	}
	return res
}
