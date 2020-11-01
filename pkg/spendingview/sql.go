package spendingview

import (
	"fmt"

	"github.com/bbokorney/budget-api-http/pkg/models"
	"github.com/bbokorney/budget-api-http/pkg/sqlutil"
	"gorm.io/gorm"
)

func QueryView(db *gorm.DB) (*SpendingView, error) {
	view := &SpendingView{
		Current: map[string]float64{},
		Limits:  map[string]float64{},
		Annual:  map[string]float64{},
	}
	type Row struct {
		Total    float64
		Category string
	}
	var queryResult []Row
	// rows, err := sqlutil.CurrentMonthWhereClause(
	result := sqlutil.CurrentMonthWhereClause(
		db.Model(&models.Transaction{}).
			Select("sum(amount) as total,category").
			Order("total desc").Group("category")).Find(&queryResult)
	// if err != nil {
	// 	return nil, err
	// }
	fmt.Println(result)
	if result.Error != nil {
		return nil, result.Error
	}
	fmt.Println(queryResult)
	for _, r := range queryResult {
		view.Current[r.Category] = r.Total
	}
	// var result []Row
	// var row Row
	// for rows.Next() {
	// 	err := db.ScanRows(rows, &row)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	result = append(result, row)
	// }
	fmt.Printf("%+v\n", view)
	return view, nil
}
