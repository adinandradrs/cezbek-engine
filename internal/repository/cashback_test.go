package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/chrisyxlee/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCashback_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	pool, tx := pgxpoolmock.NewMockPgxIface(ctrl), pgxpoolmock.NewMockPgxIface(ctrl)
	persister := NewCashback(Cashback{
		Logger: logger,
		Pool:   pool,
	})
	ctx := context.Background()
	cashback := model.Cashback{
		KezbekRefCode: sql.NullString{String: "REF001"},
		Amount:        decimal.NullDecimal{Decimal: decimal.NewFromInt(10000)},
		Reward:        decimal.NullDecimal{Decimal: decimal.NewFromInt(500)},
		WalletCode:    sql.NullString{String: "CODE_A"},
		H2HCode:       sql.NullString{String: "HOST_CODE"},
		BaseEntity: model.BaseEntity{
			CreatedBy: sql.NullInt64{Int64: 1},
		},
	}
	cmd := `INSERT INTO cashbacks 
		(kezbek_ref_code, amount, reward, wallet_code,
		h2h_code, status, is_deleted, created_by, created_date)
		VALUES ($1, $2, $3, $4, $5, $6, FALSE, $7, NOW())`
	t.Run("should success", func(t *testing.T) {
		pool.EXPECT().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}).
			Return(tx, nil)
		tx.EXPECT().Exec(ctx, cmd, cashback.KezbekRefCode.String, cashback.Amount.Decimal, cashback.Reward.Decimal,
			cashback.WalletCode.String, cashback.H2HCode.String, apps.StatusInactive,
			cashback.CreatedBy.Int64).Return(nil, nil)
		tx.EXPECT().Commit(ctx).Times(1).Return(nil)
		tx.EXPECT().Rollback(ctx).Times(1).Return(nil)
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, "failed to commit add cashback trx", r)
			}
		}()
		ex := persister.Add(cashback)
		assert.Nil(t, ex)
	})

	t.Run("should return exception on failed to begin transaction", func(t *testing.T) {
		pool.EXPECT().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}).
			Return(nil, fmt.Errorf("something went wrong"))
		ex := persister.Add(cashback)
		assert.NotNil(t, ex)
	})

	t.Run("should return exception on failed to query", func(t *testing.T) {
		pool.EXPECT().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}).
			Return(tx, nil)
		tx.EXPECT().Exec(ctx, cmd, cashback.KezbekRefCode.String, cashback.Amount.Decimal, cashback.Reward.Decimal,
			cashback.WalletCode.String, cashback.H2HCode.String, apps.StatusInactive,
			cashback.CreatedBy.Int64).Return(nil, fmt.Errorf("something went wrong"))
		tx.EXPECT().Rollback(ctx).Times(1).Return(nil)
		ex := persister.Add(cashback)
		assert.NotNil(t, ex)
	})

	t.Run("should rollback transaction on failed to commit", func(t *testing.T) {
		pool.EXPECT().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}).
			Return(tx, nil)
		tx.EXPECT().Exec(ctx, cmd, cashback.KezbekRefCode.String, cashback.Amount.Decimal, cashback.Reward.Decimal,
			cashback.WalletCode.String, cashback.H2HCode.String, apps.StatusInactive,
			cashback.CreatedBy.Int64).Return(nil, nil)
		tx.EXPECT().Commit(ctx).Times(1).Return(fmt.Errorf("something went wrong"))
		tx.EXPECT().Rollback(ctx).Times(1).Return(nil)
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, "failed to commit add cashback trx", r)
			}
		}()
		ex := persister.Add(cashback)
		assert.NotNil(t, ex)
	})
}
