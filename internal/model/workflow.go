package model

import (
	"database/sql"
	"github.com/shopspring/decimal"
)

type (
	Cashback struct {
		Id                 int64               `json:"id" db:"id"`
		MinQty             sql.NullInt32       `json:"min_qty" db:"min_qty"`
		MinTransaction     decimal.NullDecimal `json:"min_transaction" db:"min_transaction"`
		MaxTransaction     decimal.NullDecimal `json:"max_transaction" db:"max_transaction"`
		CashbackPercentage decimal.NullDecimal `json:"cashback_percentage" db:"cashback_percentage"`
		Status             int                 `json:"status" db:"status"`
		BaseEntity
	}

	Reward struct {
		Id        int64               `json:"id" db:"id"`
		Tier      sql.NullString      `json:"tier" db:"tier"`
		Reward    decimal.NullDecimal `json:"reward" db:"reward"`
		Recurring sql.NullInt32       `json:"recurring" db:"recurring"`
		Status    int                 `json:"status" db:"status"`
		BaseEntity
	}
)
