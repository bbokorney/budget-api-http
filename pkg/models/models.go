package models

import (
	"time"

	"gorm.io/gorm"
)

func AllModels() []interface{} {
	return []interface{}{
		&Transaction{},
		&Category{},
	}
}

type Transaction struct {
	gorm.Model
	Date     *time.Time
	Amount   float32
	Category string
	Vendor   string
}

type Category struct {
	gorm.Model
	Name string
}
