package handler

import (
	"confunding/helper"
	"confunding/user"
	"fmt"
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

func (h *userHandler) Login(c *gin.Context){

	var input user.LoginUserInput
	
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Register Account failed", http.StatusUnprocessableEntity, "error", errorMessage)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedinUser, err := h.service.Login(input)

	
	if err != nil {
		
		errorMessage := gin.H{"errors": err.Error()}
		
		response := helper.APIResponse("Login Account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	formatUser := user.FormatterUser(loggedinUser, "tokentokentoken")


	response := helper.APIResponse("Successfuly Login", http.StatusOK, "success", formatUser)

	c.JSON(http.StatusOK, response)

	
}

func (h *userHandler) CheckEmailAvailability(c *gin.Context){

	var user user.CheckUserInput

	err := c.ShouldBindJSON(&user)

	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.service.IsEmailAvalable(user)

	if err != nil {
		errorMessage := gin.H{"errors": "Server Error"}
		
		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{
		"is_available" : isEmailAvailable,
	}

	metaMessage := "Email has been registered"

	if isEmailAvailable {
		metaMessage = "Email is Available"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)

}

func (h *userHandler) UploadAvatar(c *gin.Context){
	//input dari user 
	// simpan gambarny di folder "images"
	// di service panggil repo
	// JWT (sementara hardcode) menggunakan id = 1
	// ambil data user id = 1
	// repo updated data user simpan lokasi file

	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{
			"is_uploaded" : false,
		}
		response := helper.APIResponse("Failded to avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return 
	}
	// harusnya dapat dari jwt 
	userId := 17
	path := fmt.Sprintf("images/%d-%s", userId, file.Filename)
	
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{
			"is_uploaded" : false,
		}
		response := helper.APIResponse("Failded to avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return 
	}

	
	_, err = h.service.SaveAvatar(userId, path)
	if err != nil {
		data := gin.H{
			"is_uploaded" : false,
		}
		response := helper.APIResponse("Failded to avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return 
	}

	data := gin.H{
		"is_uploaded" : true,
	}
	response := helper.APIResponse("Avatar Successfuly Uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)

}