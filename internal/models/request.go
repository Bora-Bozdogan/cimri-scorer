package models

import (
	"errors"
	"reflect"
)

type Request struct {
	ApiKey             *string
	ProductName        *string
	ProductDescription *string
	ProductImage       *string
	StoreName          *string
	Price              *int
	Stock              *int
	PopularityScore    *int
	UrgencyScore       *int
}

func (r Request) Validate() error {
	val := reflect.ValueOf(r)

	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).IsNil() {
			return errors.New("invalid json")
		}
	}

	return nil
}
