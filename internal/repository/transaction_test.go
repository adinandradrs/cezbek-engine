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

func TestTransaction_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	pool, tx := pgxpoolmock.NewMockPgxIface(ctrl), pgxpoolmock.NewMockPgxIface(ctrl)
	ctx := context.Background()
	persister := NewTransaction(Transaction{
		Logger: logger,
		Pool:   pool,
	})
	trx := model.Transaction{
		Email:          sql.NullString{String: "someone@email.net", Valid: true},
		Msisdn:         sql.NullString{String: "62812345679", Valid: true},
		Partner:        sql.NullString{String: "PT. Something Good", Valid: true},
		Amount:         decimal.New(250000, 10),
		KezbekRefCode:  sql.NullString{String: "KEZBEK/001/002/003", Valid: true},
		Qty:            20,
		PartnerRefCode: sql.NullString{String: "PARRTNER/001/002/003", Valid: true},
		WalletCode:     sql.NullString{String: "LSAJA", Valid: true},
		PartnerId:      1,
	}
	t.Run("should success", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"ID"}).AddRow(1).ToPgxRows()
		pool.EXPECT().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}).Return(tx, nil)
		tx.EXPECT().Commit(ctx).Times(1).Return(nil)
		tx.EXPECT().Rollback(ctx).Times(1).Return(nil)
		tx.EXPECT().QueryRow(context.Background(), `INSERT INTO transactions 
		(status, partner_id, partner, wallet_code, msisdn, email,
		qty, amount, partner_ref_code, kezbek_ref_code,
		is_deleted, created_by, created_date)
		VALUES ($1, $2, $3, $4, $5, $6, 
		$7, $8, $9, $10,
		FALSE, $11, NOW()) RETURNING ID`,
			apps.StatusInactive, trx.PartnerId, trx.Partner.String, trx.WalletCode.String, trx.Msisdn.String, trx.Email.String,
			trx.Qty, trx.Amount, trx.PartnerRefCode.String, trx.KezbekRefCode.String,
			trx.CreatedBy.Int64).Return(rows)
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, "failed to commit add kezbek trx", r)
			}
		}()
		tid, ex := persister.Add(trx)
		assert.Nil(t, ex)
		assert.NotNil(t, tid)
	})

	t.Run("should return exception on failed to begin transaction", func(t *testing.T) {
		pool.EXPECT().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}).
			Return(nil, fmt.Errorf("something went wrong on begin transaction"))
		tid, ex := persister.Add(trx)
		assert.Equal(t, "something went wrong on begin transaction", ex.Exception)
		assert.Nil(t, tid)
	})

	t.Run("should rollback on commit failure", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"ID"}).AddRow(1).ToPgxRows()
		pool.EXPECT().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}).Return(tx, nil)
		tx.EXPECT().QueryRow(context.Background(), `INSERT INTO transactions 
		(status, partner_id, partner, wallet_code, msisdn, email,
		qty, amount, partner_ref_code, kezbek_ref_code,
		is_deleted, created_by, created_date)
		VALUES ($1, $2, $3, $4, $5, $6, 
		$7, $8, $9, $10,
		FALSE, $11, NOW()) RETURNING ID`,
			apps.StatusInactive, trx.PartnerId, trx.Partner.String, trx.WalletCode.String, trx.Msisdn.String, trx.Email.String,
			trx.Qty, trx.Amount, trx.PartnerRefCode.String, trx.KezbekRefCode.String,
			trx.CreatedBy.Int64).Return(rows)
		tx.EXPECT().Rollback(ctx).Times(1).Return(nil)
		tx.EXPECT().Commit(ctx).Times(1).Return(fmt.Errorf("something went wrong on commit insert trx tx"))
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, "failed to commit add kezbek trx", r)
			}
		}()
		tid, ex := persister.Add(trx)
		assert.NotNil(t, ex)
		assert.Nil(t, tid)
	})
}
