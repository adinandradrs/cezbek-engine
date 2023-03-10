package model

import (
	"database/sql"
)

type BaseEntity struct {
	IsDeleted   bool          `json:"is_deleted"`
	CreatedBy   sql.NullInt64 `json:"created_by"`
	CreatedDate sql.NullTime  `json:"created_date"`
	UpdatedBy   sql.NullInt64 `json:"updated_by"`
	UpdatedDate sql.NullTime  `json:"updated_date"`
}

type (
	BusinessError struct {
		ErrorCode    string
		ErrorMessage string
	}

	TechnicalError struct {
		Exception string `json:"exception"`
		Occurred  int64  `json:"occurred_unixts"`
		Ticket    string `json:"ticket"`
	}
)
