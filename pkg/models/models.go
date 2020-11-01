package models

import (
	"time"
)

func AllModels() []interface{} {
	return []interface{}{
		&Transaction{},
		&Category{},
		&CategoryLimit{},
	}
}

type Transaction struct {
	ID       uint      `gorm:"primarykey" json:"id"`
	Date     time.Time `json:"date"`
	Amount   float32   `json:"amount"`
	Category string    `json:"category"`
	Vendor   string    `json:"vendor"`
}

type Category struct {
	ID   uint   `gorm:"primarykey" json:"id"`
	Name string `gorm:"unique" json:"name"`
}

type CategoryLimit struct {
	ID    uint    `gorm:"primarykey" json:"id"`
	Name  string  `gorm:"unique" json:"name"`
	Limit float64 `json:"limit"`
}
