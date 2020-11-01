package sqlutil

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

func CurrentMonthWhereClause(db *gorm.DB) *gorm.DB {
	today := time.Now()
	year, month := today.Year(), today.Month()
	start := fmt.Sprintf("%d-%d-01 00:00", year, month)
	end := fmt.Sprintf("%d-%d-01 00:00", year, (month+1)%12)
	return db.Where("date BETWEEN ? AND ?", start, end)
}

func CurrentYearWhereClause(db *gorm.DB) *gorm.DB {
	year := time.Now().Year()
	start := fmt.Sprintf("%d-01-01 00:00", year)
	end := fmt.Sprintf("%d-12-31 00:00", year)
	return db.Where("date BETWEEN ? AND ?", start, end)
}
