package handler

import (
	"confunding/campaign"
	"confunding/helper"
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

	formatter := campaign.FormatCampaignDetail(campaignDetail)

	response := helper.APIResponse("campaign detail", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}