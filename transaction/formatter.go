package transaction

import (
	"time"
)

type CampaignTransactionFormatter struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Amount int `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func FormatCampaignTransaction(transaction Transactions) CampaignTransactionFormatter{
	formatter := CampaignTransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.Name = transaction.User.Name
	formatter.Amount = transaction.Amount
	formatter.CreatedAt = transaction.CreatedAt
	
	return formatter
	
}

func FormatCampaignsTransactions(transactions []Transactions) []CampaignTransactionFormatter{
	if len(transactions) == 0 {
		return []CampaignTransactionFormatter{}
	}

	var transactionsFormatter []CampaignTransactionFormatter

	for _, transaction := range transactions {
		formatter := FormatCampaignTransaction(transaction)
		transactionsFormatter = append(transactionsFormatter, formatter)
	}

	return transactionsFormatter
}