package model

import "database/sql"

type (
	Partner struct {
		Id          int64          `json:"id" db:"id"`
		Partner     sql.NullString `json:"partner" db:"partner"`
		Code        sql.NullString `json:"code" db:"code"`
		ApiKey      sql.NullString `json:"api_key" db:"api_key"`
		Salt        sql.NullString `json:"salt" db:"salt"`
		Secret      sql.NullString `json:"secret" db:"secret"`
		Email       sql.NullString `json:"email" db:"email"`
		Msisdn      sql.NullString `json:"msisdn" db:"msisdn"`
		Officer     sql.NullString `json:"officer" db:"officer"`
		Address     sql.NullString `json:"address" db:"address"`
		PartnerLogo sql.NullString `json:"partner_logo" db:"partner_logo"`
		Status      int            `json:"status" db:"status"`
		BaseEntity
	}
)

type (
	PartnerRequest struct {
		Partner string `json:"partner" validate:"required"`
		Email   string `json:"email" validate:"required"`
		Msisdn  string `json:"msisdn" validate:"required"`
		Officer string `json:"officer" validate:"required"`
		Address string `json:"address" validate:"required"`
		Logo    string `json:"logo" validate:"required"`
	}

	PartnerAuthenticationRequest struct {
		Signature string `swaggerignore:"true"`
		Code      string `json:"code"`
	}

	PartnerOfficerAuthValidationRequest struct {
		Email string `json:"email"`
	}

	PartnerOfficerAuthVerificationRequest struct {
		Email string `json:"email"`
		Otp   string `json:"otp"`
	}
)
