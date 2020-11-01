package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bbokorney/budget-api-http/pkg/models"
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
	db     *gorm.DB
	logger *zap.Logger
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
		today := time.Now()
		year, month := today.Year(), today.Month()
		start := fmt.Sprintf("%d-%d-01 00:00", year, month)
		end := fmt.Sprintf("%d-%d-01 00:00", year, (month+1)%12)
		bs.logger.Debug("ListTransactions", zap.Any("start", start), zap.Any("end", end))
		result = bs.db.Where("date BETWEEN ? AND ?", start, end).Find(&t)
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
