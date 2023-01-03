package model

import (
	"database/sql"
	"mime/multipart"
)

type (
	Partner struct {
		Id      int64          `json:"id" db:"id"`
		Partner sql.NullString `json:"partner" db:"partner"`
		Code    sql.NullString `json:"code" db:"code"`
		ApiKey  sql.NullString `json:"api_key" db:"api_key"`
		Salt    sql.NullString `json:"salt" db:"salt"`
		Secret  sql.RawBytes   `json:"secret" db:"secret"`
		Email   sql.NullString `json:"email" db:"email"`
		Msisdn  sql.NullString `json:"msisdn" db:"msisdn"`
		Officer sql.NullString `json:"officer" db:"officer"`
		Address sql.NullString `json:"address" db:"address"`
		Logo    sql.NullString `json:"logo" db:"logo"`
		Status  int            `json:"status" db:"status"`
		BaseEntity
	}
)

type (
	AddPartnerRequest struct {
		Partner string               `json:"partner" validate:"required"`
		Code    string               `json:"code" validate:"required"`
		Email   string               `json:"email" validate:"required"`
		Msisdn  string               `json:"msisdn" validate:"required"`
		Officer string               `json:"officer" validate:"required"`
		Address string               `json:"address" validate:"required"`
		Logo    multipart.FileHeader `swaggerignore:"true" validate:"required"`
		SessionRequest
	}

	UpdatePartnerRequest struct {
		Id      int64
		Partner string                `json:"partner" validate:"required"`
		Msisdn  string                `json:"msisdn" validate:"required"`
		Officer string                `json:"officer" validate:"required"`
		Address string                `json:"address" validate:"required"`
		Logo    *multipart.FileHeader `swaggerignore:"true"`
		SessionRequest
	}

	ClientAuthenticationRequest struct {
		ApiKey string `swaggerignore:"true"`
		Code   string `json:"code" example:"LAJADA" validate:"required"`
	}

	OfficerAuthenticationRequest struct {
		Email string `json:"email" example:"john.doe@lajada.id"`
	}

	OfficerValidationRequest struct {
		TransactionId string `swaggerignore:"true" validate:"required"`
		Otp           string `json:"otp" example:"123456" validate:"required"`
	}
)

type (
	ClientAuthenticationResponse struct {
		Id      *int64 `json:"id,omitempty"`
		Code    string `json:"code"`
		Company string `json:"company"`
		SessionResponse
	}

	OfficerAuthenticationResponse struct {
		RemainingSeconds float64 `json:"remaining_seconds"`
		TransactionResponse
	}

	OfficerValidationResponse struct {
		UrlLogo string `json:"url_logo"`
		Msisdn  string `json:"msisdn"`
		Email   string `json:"email"`
		Code    string `json:"code"`
		Company string `json:"company"`
		SessionResponse
	}
)
