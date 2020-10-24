package server

import (
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

}
