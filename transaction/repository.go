package transaction

import "gorm.io/gorm"

type Repository interface {
	GetByCampaingID(campaignID int)([]Transactions, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository{
	return &repository{db}
}

func(r *repository) GetByCampaingID(campaignID int)([]Transactions, error){
	var transaction []Transactions
	err := r.db.Preload("User").Where("campaign_id = ?",campaignID).Order("id desc").Find(&transaction).Error
	if err != nil {
		return transaction, err
	}
	return transaction, nil 
}

