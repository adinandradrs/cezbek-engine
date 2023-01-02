package model

import (
	"database/sql"
	"github.com/shopspring/decimal"
)

type (
	H2HProvider struct {
		Id       int64          `json:"id" db:"id"`
		Provider sql.NullString `json:"provider" db:"provider"`
		Code     sql.NullString `json:"code" db:"code"`
		Status   int
		BaseEntity
	}

	H2HProviderFee struct {
		Id            int64               `json:"id" db:"id"`
		H2HProviderId int64               `json:"h2h_provider_id" db:"h2h_provider_id"`
		WalletCode    sql.NullString      `json:"wallet_code" db:"wallet_code"`
		Fee           decimal.NullDecimal `json:"fee" db:"fee"`
		BaseEntity
	}

	H2HPricingProjection struct {
		Code       string          `json:"code" db:"code"`
		Provider   string          `json:"provider" db:"provider"`
		WalletCode string          `json:"wallet_code" db:"wallet_code"`
		Fee        decimal.Decimal `json:"fee" db:"fee"`
	}

	H2HPricingsProjection struct {
		WalletCode string `json:"wallet_code,omitempty"`
		Prices     []H2HPricingProjection
	}
)

type (
	H2HProviderResponse struct {
		Id       int64  `json:"id,omitempty"`
		Code     string `json:"code,omitempty"`
		Provider string `json:"provider,omitempty"`
	}
)
