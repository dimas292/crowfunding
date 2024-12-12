package handler

import (
	"confunding/campaign"
	"confunding/helper"
	"confunding/user"
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

func(h *campaignHandler) GetCampaign(c *gin.Context){
	// hanlder menangkap id dan membuat response formatter 
	// service menerima id untuk dilanjutkan ke repo
	// repository membuat find by id
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

	response := helper.APIResponse("campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}

func(h *campaignHandler) CreateCampaign(c *gin.Context){
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

func(h *campaignHandler) UpdateCampaign(c *gin.Context){
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
