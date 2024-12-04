package handler

import (
	"confunding/helper"
	"confunding/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	service user.Service
}

func NewUserHanlder (userService user.Service) *userHandler{
	return &userHandler{userService}
}

func (h userHandler) RegisterUser(c *gin.Context){
	// tangkap input dari user
	// map input dari user ke struct registerinput
	// struct di atas kita passing ke parameter service

	var input user.RegisterUserInput
	
	err := c.ShouldBindJSON(&input)


	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Register Account failed", http.StatusUnprocessableEntity, "error", errorMessage)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	newUser, err := h.service.RegisterUser(input)

	userFormatter := user.FormatterUser(newUser, "tokentokentokentoken")



	if err != nil {
		response := helper.APIResponse("Register Account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	
	response := helper.APIResponse("Account has been register!", http.StatusOK, "success", userFormatter)

	c.JSON(http.StatusOK, response)
}