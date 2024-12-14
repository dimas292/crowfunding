package main

import (
	"confunding/auth"
	"confunding/campaign"
	"confunding/handler"
	"confunding/helper"
	"confunding/user"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	dsn := "root:@tcp(127.0.0.1:3306)/confunding?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	authService := auth.NewService()
	userHandler := handler.NewUserHanlder(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	router := gin.Default()
	router.Static("/images", "./images")
	
	api := router.Group("/api/v1")
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService),userHandler.UploadAvatar)
	// campaigns
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaign/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService),campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadCampaignImage)
	router.Run()

}

// wraping middleware 
func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc{
	return func (c *gin.Context){
		authHeader := c.GetHeader("Authorization")
		// validasi apakah didalam auth header ada string bernama bearer
		if !strings.Contains(authHeader, "Bearer"){
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
	
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		tokenString := ""
		arrayString := strings.Split(authHeader, " ")
		if len(arrayString) == 2 {
			tokenString = arrayString[1]
		}
		
		// memvalidasi token dengan secret key 
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
	
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// merubah token kewujud aseli

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid{
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
	
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		//mendapatkan nilai user id dan parsing ke float
		userId := int(claim["user_id"].(float64))
		// mendapatkan nilai user dengan find by id
		user, err := userService.GetUserById(userId)

		if err != nil{
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
	
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		// set context dengan key currentUser dan valuenya user
		c.Set("currentUser", user)

	}


}

