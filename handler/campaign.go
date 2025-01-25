package handler

import (
	"confunding/campaign"
	"confunding/helper"
	"confunding/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// tangkap parameter di hanler
// service menentukan repository yang dipakai
// repository
// db

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := campaign.FormatCampaigs(campaigns)

	response := helper.APIResponse("list of campaigns", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)

}

// campaign details

func (h *campaignHandler) GetCampaign(c *gin.Context) {

	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Error to get campaigns id", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err := h.service.GetCampaignByID(input)
	if err != nil {
		response := helper.APIResponse("Error to get campaigns id", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	if campaignDetail.ID == 0 {
		response := helper.APIResponse("Campaign not found", http.StatusNotFound, "error", nil)
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := helper.APIResponse("campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity, "error", errorMessage)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser

	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity, "error", nil)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.APIResponse("Success to create campaign", http.StatusOK, "success", campaign.FormatCampaign(newCampaign))

	c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	var inputID campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("Fail to update campaign", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData campaign.CreateCampaignInput

	err = c.ShouldBind(&inputData)

	if err != nil {
		response := helper.APIResponse("Fail to update campaign", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}
	currentUser := c.MustGet("currentUser").(user.User)
	inputData.User = currentUser

	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		response := helper.APIResponse("Fail to update campaign", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Succes Update Campaign", http.StatusOK, "success", campaign.FormatCampaign(updatedCampaign))

	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UploadCampaignImage(c *gin.Context) {
	var input campaign.CreateCampaignImageInput

	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity, "error", errorMessage)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helper.APIResponse("Failded upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	userId := currentUser.ID
	path := fmt.Sprintf("images/%d-%s", userId, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{
			"is_uploaded" : false,
		}
		response := helper.APIResponse("Failded to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return 
	}
	input.User = currentUser
	_, err = h.service.SaveCampaignImages(input, path)
	if err != nil {
		data := gin.H{
			"is_uploaded" : false,
		}
		response := helper.APIResponse("Failded to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return 
	}

	data := gin.H{
		"is_uploaded" : true,
	}
	response := helper.APIResponse("Image Successfuly Uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
