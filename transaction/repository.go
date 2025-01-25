package transaction

import "gorm.io/gorm"

type Repository interface {
	GetByCampaingID(campaignID int)([]Transactions, error)
	GetByUserID(userID int)([]Transactions, error)
	Save(transaction Transactions)(Transactions, error)
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

func (r *repository) GetByUserID(userID int)([]Transactions, error){
	var transactions []Transactions
	err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userID).Order("id desc").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil 
}

func(r *repository) Save(transaction Transactions)(Transactions, error) {
	err := r.db.Create(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil 
}



