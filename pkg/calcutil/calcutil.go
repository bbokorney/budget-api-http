package calcutil

import "time"

func UnplannedMonthlySpending(annualLimit, currentYearSpending, annualPlannedSpending, monthlyPlannedSpeding float64) float64 {
	return (annualLimit-currentYearSpending-annualPlannedSpending)/float64(12-time.Now().Month()+1) - monthlyPlannedSpeding
}
