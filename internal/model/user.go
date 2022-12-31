package model

import "database/sql"

type (
	User struct {
		Id       int64          `json:"id" db:"id"`
		Fullname sql.NullString `json:"fullname" db:"fullname"`
		Msisdn   sql.NullString `json:"msisdn" db:"msisdn"`
		Email    sql.NullString `json:"email" db:"email"`
		SubId    sql.NullString `json:"sub_id" db:"sub_id"`
		Status   int            `json:"status" db:"status"`
		BaseEntity
	}
)

type (
	AddUserRequest struct {
		Fullname string `json:"fullname" validate:"required"`
		Msisdn   string `json:"msisdn" validate:"required"`
		Email    string `json:"email" validate:"required"`
		RoleId   int64  `json:"role_id" validate:"required"`
		SessionRequest
	}

	UpdateUserRequest struct {
		Id       int64  `json:"id" validate:"required"`
		Fullname string `json:"fullname" validate:"required"`
		Msisdn   string `json:"msisdn" validate:"required"`
		RoleId   int64  `json:"role_id" validate:"required"`
		SessionRequest
	}
)
