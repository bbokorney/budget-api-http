package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bbokorney/budget-api-http/pkg/models"
	"github.com/bbokorney/budget-api-http/pkg/server"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	logger := buildLogger()
	logger.Info("Setting up database")
	dbDirectory := "./dbs"
	dbFilePath := filepath.Join(dbDirectory, "budget.db")
	if err := os.MkdirAll(dbDirectory, 0755); err != nil {
		logger.Fatal("Failed to create database directory", zap.Error(err))
	}

	db, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	for _, m := range models.AllModels() {
		if err := db.AutoMigrate(m); err != nil {
			logger.Fatal("Failed to auto-migrate schema", zap.Error(err))
		}
	}

	logger.Info("Database setup complete")

	listenAddr := ":8000"
	logger.Info("Starting listener", zap.Any("address", listenAddr))

	bs := server.NewBudgetServer(db, logger)

	r := gin.Default()
	r.POST("/v1/transactions", bs.AddTransaction)

	if err := r.Run(listenAddr); err != nil {
		logger.Fatal("Error running server", zap.Error(err))
	}
}

func buildLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.Level.SetLevel(zapcore.InfoLevel)

	var err error
	logger, err := config.Build()
	if err != nil {
		fmt.Printf("Error setting up logger: %s\n", err)
		os.Exit(1)
	}

	return logger
}
