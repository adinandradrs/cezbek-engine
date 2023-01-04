package model

import (
	"database/sql"
	"github.com/shopspring/decimal"
)

type (
	Cashback struct {
		Id            int64               `json:"id" db:"id"`
		KezbekRefCode sql.NullString      `json:"kezbek_ref_code" db:"kezbek_ref_code"`
		Amount        decimal.NullDecimal `json:"amount" db:"amount"`
		Reward        decimal.NullDecimal `json:"reward" db:"reward"`
		WalletCode    sql.NullString      `json:"wallet_code" db:"wallet_code"`
		H2HCode       sql.NullString      `json:"h2h_code" db:"h2h_code"`
		BaseEntity
	}

	Tier struct {
		Id                   int64          `json:"id" db:"id"`
		PartnerId            int64          `json:"partner_id" db:"partner_id"`
		Msisdn               sql.NullString `json:"msisdn" db:"msisdn"`
		Email                sql.NullString `json:"email" db:"email"`
		NextGrade            int            `json:"next_grade" db:"next_grade"`
		NextTier             sql.NullString `json:"next_tier" db:"next_tier"`
		CurrentGrade         int            `json:"current_grade" db:"current_grade"`
		CurrentTier          sql.NullString `json:"current_tier" db:"current_tier"`
		PrevGrade            int            `json:"prev_grade" db:"prev_grade"`
		PrevTier             sql.NullString `json:"prev_tier" db:"prev_tier"`
		ExpiredDate          sql.NullTime   `json:"expired_date" db:"expired_date"`
		TransactionRecurring int            `json:"transaction_recurring" db:"transaction_recurring"`
		Journey              TierJourneys
		BaseEntity
	}

	TierJourneys struct {
		Id                int64          `json:"id" db:"id"`
		CurrentGrade      int            `json:"current_grade" db:"current_grade"`
		CurrentTier       sql.NullString `json:"current_tier" db:"current_tier"`
		LastTransactionId int64          `json:"last_transaction_id" db:"last_transaction_id"`
		Notes             sql.NullString `json:"notes" db:"notes"`
		TierId            int64          `json:"tier_id" db:"tier_id"`
		BaseEntity
	}

	WfRewardTierGradeProjection struct {
		Tier  *string `json:"tier,omitempty" db:"tier"`
		Grade *int    `json:"grade,omitempty" db:"grade"`
	}

	WfRewardTierProjection struct {
		Recurring    int                         `json:"recurring,omitempty" db:"recurring"`
		MaxRecurring int                         `json:"max_recurring,omitempty" db:"max_recurring"`
		Tier         string                      `json:"tier,omitempty" db:"tier"`
		Grade        int                         `json:"grade,omitempty" db:"grade"`
		PrevTier     WfRewardTierGradeProjection `json:"prev_tier,omitempty" db:"prev_tier"`
		NextTier     WfRewardTierGradeProjection `json:"next_tier,omitempty" db:"next_tier"`
		Reward       decimal.Decimal             `json:"reward,omitempty" db:"reward"`
	}
)

type (
	TierRequest struct {
		PartnerId     int64
		Msisdn        string
		Email         string
		TransactionId int64
	}

	FindCashbackRequest struct {
		Qty    int
		Amount decimal.Decimal
	}
)

type (
	FindCashbackResponse struct {
		Amount decimal.Decimal
	}
)
