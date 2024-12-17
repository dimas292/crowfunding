package transaction

import (
	"confunding/campaign"
	"errors"
)

type Service interface {
	GetTransactionByCampaignID(input GetCampaignTransactionsInput)([]Transactions, error)
	GetTransactionByUserID(userID int)([]Transactions, error)
}

type service struct {
	repository Repository
	campaignRepository campaign.Repository
}

func NewService(repository Repository, campaignRepository campaign.Repository) *service{
	return &service{repository, campaignRepository}
}

func(s *service) GetTransactionByCampaignID(input GetCampaignTransactionsInput)([]Transactions, error){

	campaign, err  := s.campaignRepository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}


	if campaign.UserID != input.User.ID {
		return []Transactions{}, errors.New("not an ownwer of the campaign")

	}
	
	transactions, err := s.repository.GetByCampaingID(input.ID)
	if err != nil {
		return transactions, err
	}
	
	return transactions, nil 
}

func(s *service) GetTransactionByUserID(userID int)([]Transactions, error) {
	
	transactions, err := s.repository.GetByUserID(userID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil 
}
