package model

import (
	"database/sql"
	"github.com/shopspring/decimal"
)

type (
	Transaction struct {
		Id             int64           `json:"id" db:"id"`
		Status         int             `json:"status" db:"status"`
		PartnerId      int64           `json:"partner_id" db:"partner_id"`
		Partner        sql.NullString  `json:"partner" db:"partner"`
		WalletCode     sql.NullString  `json:"wallet_code" db:"wallet_code"`
		Msisdn         sql.NullString  `json:"msisdn" db:"msisdn"`
		Email          sql.NullString  `json:"email" db:"email"`
		Qty            int             `json:"qty" db:"qty"`
		Amount         decimal.Decimal `json:"amount" db:"amount"`
		PartnerRefCode sql.NullString  `json:"partner_ref_code" db:"kezbek_ref_code"`
		KezbekRefCode  sql.NullString  `json:"kezbek_ref_code" db:"kezbek_ref_code"`
		BaseEntity
	}

	PartnerTransactionProjection struct {
		Id          int64           `json:"id" example:"1"`
		WalletCode  string          `json:"wallet_code,omitempty" example:"LSAJA"`
		Email       string          `json:"email,omitempty" example:"john.doe@email.net"`
		Msisdn      string          `json:"msisdn,omitempty" example:"628118770510"`
		Qty         int             `json:"qty,omitempty" example:"2"`
		Transaction decimal.Decimal `json:"transaction,omitempty" example:"250000"`
		Cashback    decimal.Decimal `json:"cashback,omitempty" example:"2500"`
		Reward      decimal.Decimal `json:"reward,omitempty" example:"13000"`
	}
)

type (
	PartnerTransactionSearchResponse struct {
		Transactions []PartnerTransactionProjection `json:"transactions,omitempty"`
		PaginationResponse
	}

	TransactionTierResponse struct {
		Tier        string `json:"tier,omitempty" example:"GOLD"`
		Recurring   int    `json:"recurring,omitempty" example:"3"`
		DateExpired string `json:"date_expired,omitempty" example:"2022-01-01"`
	}
)

type (
	TransactionRequest struct {
		Qty                  int             `json:"quantity" example:"2" validate:"required"`
		Amount               decimal.Decimal `json:"amount" example:"750000" validate:"required"`
		Msisdn               string          `json:"msisdn" example:"62812345678" validate:"required"`
		Email                string          `json:"email" example:"john.doe@gmailxyz.com"`
		MerchantCode         string          `json:"merchant_code" example:"LSAJA,GPAID,JOSVO"`
		TransactionReference string          `json:"transaction_reference" example:"INV/001/002"`
		SessionRequest
	}
)
