package transaction

import (
	"confunding/user"
	"time"
)

type Transactions struct {
	ID int
	CampaignID int
	UserID int
	Amount int
	Status string
	Code string
	User user.User
	CreatedAt time.Time
	UpdatedAt time.Time
}
