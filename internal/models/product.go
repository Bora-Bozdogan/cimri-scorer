package models

type Product struct {
	ID int 
	ProductName string `gorm:"column:name"`
	ProductDescription string `gorm:"column:description"`
	ProductImage string `gorm:"column:image_url"`
	Popularity_score int `gorm:"column:popularity_score"`
	Urgency_score int `gorm:"column:urgency_score"`
}
