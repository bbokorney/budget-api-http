package server

import (
	"fmt"
	"net/http"

	"github.com/bbokorney/budget-api-http/pkg/calcutil"
	"github.com/bbokorney/budget-api-http/pkg/models"
	"github.com/bbokorney/budget-api-http/pkg/spendingview"
	"github.com/bbokorney/budget-api-http/pkg/sqlutil"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewBudgetServer(db *gorm.DB, logger *zap.Logger) *BudgetServer {
	return &BudgetServer{
		db:     db,
		logger: logger,
	}
}

type BudgetServer struct {
	db                *gorm.DB
	logger            *zap.Logger
	spendingViewCache *spendingview.Container
}

func (bs *BudgetServer) AddTransaction(c *gin.Context) {
	var t models.Transaction
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bs.logger.Debug("AddTransaction", zap.Any("transaction", t))
	result := bs.db.Create(&t)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusAccepted)
}

func (bs *BudgetServer) ListTransactions(c *gin.Context) {
	type listTransactionsParams struct {
		CurrentMonth bool `form:"current_month"`
	}
	var params listTransactionsParams
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bs.logger.Debug("ListTransactions params", zap.Any("params", params))
	var t []models.Transaction
	var result *gorm.DB
	if params.CurrentMonth {
		result = bs.db.Scopes(sqlutil.CurrentMonthWhereClause).Find(&t)
	} else {
		result = bs.db.Find(&t)
	}
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	bs.logger.Debug("ListTransactions", zap.Any("transactions", t))
	c.JSON(http.StatusOK, t)
}

func (bs *BudgetServer) AddCategory(c *gin.Context) {
	var cat models.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bs.logger.Debug("AddCategory", zap.Any("category", cat))
	result := bs.db.Create(&cat)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusAccepted)
}

func (bs *BudgetServer) ListCategories(c *gin.Context) {
	var cats []models.Category
	result := bs.db.Find(&cats)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	bs.logger.Debug("ListCategory", zap.Any("categories", cats))
	c.JSON(http.StatusOK, cats)
}

func (bs *BudgetServer) GetSpending(c *gin.Context) {
	type Row struct {
		Total    float64 `json:"amount"`
		Category string  `json:"name"`
	}
	var queryResult []Row
	result := bs.db.Model(&models.Transaction{}).
		Scopes(sqlutil.CurrentMonthWhereClause).
		Select("sum(amount) as total,category").
		Group("category").
		Order("total desc").
		Find(&queryResult)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	retBody := map[string]float64{}
	totalSum := 0.0
	for _, r := range queryResult {
		retBody[r.Category] = r.Total
		totalSum += r.Total
	}
	retBody["Total"] = totalSum
	c.JSON(http.StatusOK, retBody)
}

func (bs *BudgetServer) AddCategoryLimit(c *gin.Context) {
	var limit models.CategoryLimit
	if err := c.ShouldBindJSON(&limit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bs.logger.Debug("AddCategoryLimit", zap.Any("limit", limit))
	result := bs.db.Create(&limit)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusAccepted)
}

func (bs *BudgetServer) ListCategoryLimits(c *gin.Context) {
	var limits []models.CategoryLimit
	result := bs.db.Find(&limits)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	bs.logger.Debug("ListCategoryLimits", zap.Any("limits", limits))

	// TODO: dry up this code
	// TODO: make one endpoint of these three spending related endpoints
	var annualLimit float64
	result = bs.db.Model(&models.AnnualLimit{}).
		Select("min(amount)").First(&annualLimit)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	bs.logger.Debug("ListCategoryLimits", zap.Any("annualLimit", annualLimit))

	var annualTotal float64
	result = bs.db.Model(&models.Transaction{}).
		Scopes(sqlutil.CurrentYearWhereClause).
		Select("sum(amount)").
		Find(&annualTotal)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	bs.logger.Debug("ListCategoryLimits", zap.Any("annualTotal", annualTotal))

	var annualPlannedSpending float64
	result = bs.db.Model(&models.AnnualPlannedSpending{}).
		Select("min(amount)").First(&annualPlannedSpending)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	bs.logger.Debug("ListCategoryLimits", zap.Any("annualPlannedSpending", annualPlannedSpending))

	var monthlyPlannedTotal float64
	result = bs.db.Model(&models.CategoryLimit{}).
		Select("sum(amount)").
		Find(&monthlyPlannedTotal)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	bs.logger.Debug("ListCategoryLimits", zap.Any("monthlyPlannedTotal", monthlyPlannedTotal))

	var previousMonthsTotalSpending float64
	result = bs.db.Model(&models.Transaction{}).
		Scopes(sqlutil.AllPreviousMonthsWhereClause).
		Select("sum(amount)").
		Find(&previousMonthsTotalSpending)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	bs.logger.Debug("ListCategoryLimits", zap.Any("previousMonthsTotalSpending", previousMonthsTotalSpending))

	retBody := map[string]float64{}
	otherSpending := calcutil.UnplannedMonthlySpending(annualLimit, previousMonthsTotalSpending, annualPlannedSpending, monthlyPlannedTotal)

	retBody["Other"] = otherSpending
	var total float64 = otherSpending
	for _, l := range limits {
		retBody[l.Name] = l.Amount
		total += l.Amount
	}
	retBody["Total"] = total
	c.JSON(http.StatusOK, retBody)
}

func (bs *BudgetServer) AddAnnualLimit(c *gin.Context) {
	var limit models.AnnualLimit
	if err := c.ShouldBindJSON(&limit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bs.logger.Debug("AddAnnualLimit", zap.Any("limit", limit))
	result := bs.db.Create(&limit)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusAccepted)
}

func (bs *BudgetServer) ListAnnualLimits(c *gin.Context) {
	var limits []models.AnnualLimit
	result := bs.db.Find(&limits)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	bs.logger.Debug("ListAnnualLimits", zap.Any("limits", limits))

	var annualTotal float64
	result = bs.db.Model(&models.Transaction{}).
		Scopes(sqlutil.CurrentYearWhereClause).
		Select("sum(amount)").
		Find(&annualTotal)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	bs.logger.Debug("ListAnnualLimits", zap.Any("annualTotal", annualTotal))

	var monthlyPlannedTotal float64
	result = bs.db.Model(&models.CategoryLimit{}).
		Select("sum(amount)").
		Find(&monthlyPlannedTotal)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	bs.logger.Debug("ListAnnualLimits", zap.Any("monthlyPlannedTotal", monthlyPlannedTotal))

	var annualPlannedSpending float64
	result = bs.db.Model(&models.AnnualPlannedSpending{}).
		Select("min(amount)").First(&annualPlannedSpending)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	bs.logger.Debug("ListAnnualLimits", zap.Any("annualPlannedSpending", annualPlannedSpending))

	var previousMonthsTotalSpending float64
	result = bs.db.Model(&models.Transaction{}).
		Scopes(sqlutil.AllPreviousMonthsWhereClause).
		Select("sum(amount)").
		Find(&previousMonthsTotalSpending)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	bs.logger.Debug("ListAnnualLimits", zap.Any("previousMonthsTotalSpending", previousMonthsTotalSpending))

	retBody := map[string]float64{}
	for _, l := range limits {
		retBody[fmt.Sprintf("%.0fk", l.Amount/1000.0)] = calcutil.UnplannedMonthlySpending(l.Amount, previousMonthsTotalSpending, annualPlannedSpending, monthlyPlannedTotal)
	}
	retBody["Total"] = annualTotal
	c.JSON(http.StatusOK, retBody)
}
