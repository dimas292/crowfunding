package middleware

import (
	"confunding/auth"
	"confunding/helper"
	"confunding/user"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc{
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