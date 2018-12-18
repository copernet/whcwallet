package model

import "github.com/shopspring/decimal"

type BalanceForAddress struct {
	Address          string           `json:"address" gorm:"address"`
	PropertyID       int64            `json:"property_id" gorm:"property_id"`
	PropertyName     string           `json:"property_name" gorm:"property_name"`
	Precision        int              `json:"precision" gorm:"precision"`
	BalanceAvailable *decimal.Decimal `json:"balance_available" gorm:"balance_available"`
	Pendingpos       *decimal.Decimal `json:"pendingpos" gorm:"pendingpos"`
	Pendingneg       *decimal.Decimal `json:"pendingneg" gorm:"pendingneg"`
}
