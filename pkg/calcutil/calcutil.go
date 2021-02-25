package calcutil

import (
	"time"
)

func UnplannedMonthlySpending(annualLimit, previousMonthsTotalSpending, annualPlannedSpending, monthlyPlannedSpending float64) float64 {
	monthsLeftInYear := float64(12-time.Now().Month()) + 1
	unplannedSpendingLeft := annualLimit - previousMonthsTotalSpending - annualPlannedSpending - (monthsLeftInYear)*monthlyPlannedSpending
	return unplannedSpendingLeft / monthsLeftInYear
}
