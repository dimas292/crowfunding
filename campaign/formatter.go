package campaign

type CampaignFormatter struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	ImageUrl         string `json:"image_url"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_user"`
	Slug             string `json:"slug"`
}

func FormatCampaign(campaign Campaign) CampaignFormatter {
	campaignFormatter := CampaignFormatter{}
	campaignFormatter.ID = campaign.ID
	campaignFormatter.UserID = campaign.UserID
	campaignFormatter.Name = campaign.Name
	campaignFormatter.ShortDescription = campaign.ShortDescription
	campaignFormatter.ImageUrl = ""
	campaignFormatter.GoalAmount = campaign.GoalAmount
	campaignFormatter.Slug = campaign.Slug
	campaignFormatter.CurrentAmount = campaign.CurrentAmount

	if len(campaign.CampaignImages) > 0 {
		campaignFormatter.ImageUrl = campaign.CampaignImages[0].FileName
	}

	return campaignFormatter
}

func FormatCampaigs(campaigns []Campaign) []CampaignFormatter {

	campaignsFormatter := []CampaignFormatter{}
	for _, campaign := range campaigns {
		newCampaign := FormatCampaign(campaign)

		campaignsFormatter = append(campaignsFormatter, newCampaign)
	}

	return campaignsFormatter
}
