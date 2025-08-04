package models

import (
	"time"
)

//load all merchants with API keys into this struct
type Merchant struct {
	ID int
	Name string
	Key string
	Is_active bool 
	CreatedAt  time.Time `gorm:"column:created_at"`
}