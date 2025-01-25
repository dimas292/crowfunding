package main

import (
	"confunding/auth"
	"confunding/campaign"
	"confunding/handler"
	"confunding/middleware"
	"confunding/transaction"
	"confunding/user"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DSN string `yaml:"dsn"`
}

func main() {
	// config 
	f, err := os.ReadFile("config.yaml")
	if err != nil {
		return 
	}
	config := Config{}
	err = yaml.Unmarshal(f, &config)
	if err != nil {
		return
	}
	
	db, err := gorm.Open(postgres.Open(config.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err != nil {
		log.Fatal(err.Error())
	}
	// auto migration 

	if err := db.AutoMigrate(&campaign.Campaign{}, &campaign.CampaignImage{}, &user.User{}, &transaction.Transactions{}); err != nil {
		log.Fatalf("failed to migrate database : %v", err)
	}
	

	// repository 
	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)
	// service 
	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	authService := auth.NewService()
	transactionService := transaction.NewService(transactionRepository, campaignRepository)
	// handler
	userHandler := handler.NewUserHanlder(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	router := gin.Default()
	router.Static("/images", "./images")
	
	
	api := router.Group("/api/v1")
	// users
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", middleware.AuthMiddleware(authService, userService),userHandler.UploadAvatar)
	// campaigns
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", middleware.AuthMiddleware(authService, userService),campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", middleware.AuthMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", middleware.AuthMiddleware(authService, userService), campaignHandler.UploadCampaignImage)
	// transaction
	api.GET("/campaign/:id/transactions", middleware.AuthMiddleware(authService, userService), transactionHandler.GetCampaignTransactin)
	api.GET("/transactions", middleware.AuthMiddleware(authService, userService), transactionHandler.GetUserTransaction)
	api.POST("/transactions", middleware.AuthMiddleware(authService, userService), transactionHandler.CreateTransactions)


	router.Run()
}

