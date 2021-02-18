package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bbokorney/budget-api-http/pkg/models"
	"github.com/bbokorney/budget-api-http/pkg/server"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	authTokenHeader = "X-Auth-Token"
)

func main() {
	authToken := os.Getenv("AUTH_TOKEN")
	if authToken == "" {
		fmt.Println("Must set AUTH_TOKEN env var")
		os.Exit(1)
	}
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

	listenAddr := "127.0.0.1:8000"
	logger.Info("Starting listener", zap.Any("address", listenAddr))

	bs := server.NewBudgetServer(db, logger)

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, authTokenHeader)
	corsConfig.AddAllowMethods("OPTIONS")

	r := gin.Default()
	r.Use(cors.New(corsConfig))
	r.Use(authHandler(authToken))
	apiv1 := r.Group("/api/v1")
	{
		apiv1.POST("/transactions", bs.AddTransaction)
		apiv1.GET("/transactions", bs.ListTransactions)

		apiv1.POST("/categories", bs.AddCategory)
		apiv1.GET("/categories", bs.ListCategories)

		apiv1.POST("/category-limits", bs.AddCategoryLimit)
		apiv1.GET("/category-limits", bs.ListCategoryLimits)

		apiv1.POST("/annual-limits", bs.AddAnnualLimit)
		apiv1.GET("/annual-limits", bs.ListAnnualLimits)

		apiv1.GET("/spending", bs.GetSpending)
	}

	if err := r.Run(listenAddr); err != nil {
		logger.Fatal("Error running server", zap.Error(err))
	}
}

func authHandler(authToken string) func(c *gin.Context) {
	return func(c *gin.Context) {
		userToken := c.Request.Header.Get(authTokenHeader)
		if userToken != authToken {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

func buildLogger() *zap.Logger {
	// TODO: make logging level dynamic via rest endpoint
	// https://godoc.org/go.uber.org/zap#example-AtomicLevel
	config := zap.NewDevelopmentConfig()
	logLevel := os.Getenv("LOG_LEVEL")
	fmt.Println("loglevel =", logLevel)
	level := zapcore.DebugLevel
	switch strings.ToLower(logLevel) {
	case "debug":
		level = zapcore.DebugLevel
		break
	case "info":
		level = zapcore.InfoLevel
		break
	}
	fmt.Printf("Setting log level to %s\n", level)
	config.Level.SetLevel(level)

	var err error
	logger, err := config.Build()
	if err != nil {
		fmt.Printf("Error setting up logger: %s\n", err)
		os.Exit(1)
	}

	return logger
}
