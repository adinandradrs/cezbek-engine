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
)

type (
	TransactionRequest struct {
		Qty                  int             `json:"quantity" example:"2" validate:"required"`
		Amount               decimal.Decimal `json:"amount" example:"750000" validate:"required"`
		Msisdn               string          `json:"msisdn" example:"62812345678" validate:"required"`
		Email                string          `json:"email" example:"john.doe@gmailxyz.com"`
		WalletCode           string          `json:"wallet_code" example:"LSAJA,GPAID,JOSVO"`
		TransactionReference string          `json:"transaction_reference" example:"INV001-GOFOOD-ABC"`
	}
)
