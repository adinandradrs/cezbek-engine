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

func TestTransaction_DetailByPartner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	pool := pgxpoolmock.NewMockPgxIface(ctrl)
	ctx := context.Background()
	persister := NewTransaction(Transaction{
		Logger: logger,
		Pool:   pool,
	})
	cmd := `select t.id, t.wallet_code, t.email, t.msisdn, 
			t.qty, t.amount as transaction, c.amount as cashback, c.reward 
			from transactions t inner join cashbacks c 
			on t.kezbek_ref_code = c.kezbek_ref_code and t.id = $1 and t.partner_id = $2`
	inp := &model.FindByIdRequest{
		Id: 1,
		SessionRequest: model.SessionRequest{
			Id: 7,
		},
	}
	t.Run("should success", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"id", "wallet_code", "email", "msisdn", "qty",
			"cashback", "reward"}).AddRow(int64(1), "WALLET_A", "someone@email.net", "628118770510", 1,
			decimal.NewFromInt(5000), decimal.Zero).ToPgxRows()
		pool.EXPECT().Query(ctx, cmd, inp.Id, inp.SessionRequest.Id).Return(rows, nil)
		v, ex := persister.DetailByPartner(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})

	t.Run("should return exception on failed to query", func(t *testing.T) {
		pool.EXPECT().Query(ctx, cmd, inp.Id, inp.SessionRequest.Id).Return(nil, fmt.Errorf("something went wrong"))
		v, ex := persister.DetailByPartner(inp)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})

	t.Run("should return exception on failed to map the result", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"id", "wallet_code", "email", "msisdn", "qty",
			"cashback", "reward"}).AddRow(1, "WALLET_A", "someone@email.net", "628118770510", 1,
			decimal.NewFromInt(5000), decimal.Zero).ToPgxRows()
		pool.EXPECT().Query(ctx, cmd, inp.Id, inp.SessionRequest.Id).Return(rows, nil)
		v, ex := persister.DetailByPartner(inp)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})
}

func TestTransaction_CountByPartner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	pool := pgxpoolmock.NewMockPgxIface(ctrl)
	ctx := context.Background()
	persister := NewTransaction(Transaction{
		Logger: logger,
		Pool:   pool,
	})

	cmd := `select count(t.id) 
			from transactions t inner join cashbacks c 
			on t.kezbek_ref_code = c.kezbek_ref_code
			where t.partner_id = $1 `
	inp := &model.SearchRequest{
		SessionRequest: model.SessionRequest{
			Id: 7,
		},
	}
	t.Run("should success without text search", func(t *testing.T) {
		where := " AND '1' = $2 "
		rows := pgxpoolmock.NewRows([]string{"count"}).
			AddRow(10).ToPgxRows()
		pool.EXPECT().QueryRow(ctx, cmd+where, inp.SessionRequest.Id, gomock.Any()).
			Return(rows)
		v, ex := persister.CountByPartner(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})

	t.Run("should success with text search", func(t *testing.T) {
		inp.TextSearch = "something to search"
		where := ` AND (t.msisdn like $2 OR t.email like $2 OR t.kezbek_ref_code like $2 ) `
		rows := pgxpoolmock.NewRows([]string{"count"}).
			AddRow(10).ToPgxRows()
		pool.EXPECT().QueryRow(ctx, cmd+where, inp.SessionRequest.Id, gomock.Any()).
			Return(rows)
		v, ex := persister.CountByPartner(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})
}

func TestTransaction_SearchByPartner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	pool := pgxpoolmock.NewMockPgxIface(ctrl)
	ctx := context.Background()
	persister := NewTransaction(Transaction{
		Logger: logger,
		Pool:   pool,
	})
	inp := &model.SearchRequest{
		TextSearch: "",
		Start:      1,
		Limit:      10,
		SortBy:     "",
		Sort:       "",
		SessionRequest: model.SessionRequest{
			Id: 1,
		},
	}
	cmd := `select t.id, t.wallet_code, t.email, t.msisdn, 
			t.qty, t.amount as transaction, c.amount as cashback, c.reward 
			from transactions t inner join cashbacks c 
			on t.kezbek_ref_code = c.kezbek_ref_code
			where t.partner_id = $1  AND '1' = $2 ` + `
			order by `

	t.Run("should success without text search", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"id", "wallet_code", "email", "msisdn", "qty", "transaction",
			"cashback", "reward"}).
			AddRow(int64(1), "CODE_A", "someone@email.id", "628118770510", 1, decimal.NewFromInt(25000),
				decimal.NewFromInt(2500), decimal.NewFromInt(1000)).ToPgxRows()
		pool.EXPECT().Query(ctx, cmd+" t.id   DESC  limit $3 offset $4", inp.SessionRequest.Id, "1", inp.Limit, inp.Start).Return(rows, nil)
		v, ex := persister.SearchByPartner(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})

	t.Run("should success without text search and sorted by wallet code", func(t *testing.T) {
		inp.SortBy = "WALLET_CODE"
		rows := pgxpoolmock.NewRows([]string{"id", "wallet_code", "email", "msisdn", "qty", "transaction",
			"cashback", "reward"}).
			AddRow(int64(1), "CODE_A", "someone@email.id", "628118770510", 1, decimal.NewFromInt(25000),
				decimal.NewFromInt(2500), decimal.NewFromInt(1000)).ToPgxRows()
		pool.EXPECT().Query(ctx, cmd+"t.wallet_code  DESC  limit $3 offset $4", inp.SessionRequest.Id, "1", inp.Limit, inp.Start).Return(rows, nil)
		v, ex := persister.SearchByPartner(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})

	t.Run("should success without text search and sorted by amount", func(t *testing.T) {
		inp.SortBy = "AMOUNT"
		rows := pgxpoolmock.NewRows([]string{"id", "wallet_code", "email", "msisdn", "qty", "transaction",
			"cashback", "reward"}).
			AddRow(int64(1), "CODE_A", "someone@email.id", "628118770510", 1, decimal.NewFromInt(25000),
				decimal.NewFromInt(2500), decimal.NewFromInt(1000)).ToPgxRows()
		pool.EXPECT().Query(ctx, cmd+"c.amount  DESC  limit $3 offset $4", inp.SessionRequest.Id, "1", inp.Limit, inp.Start).Return(rows, nil)
		v, ex := persister.SearchByPartner(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})

	t.Run("should success without text search and sorted by amount", func(t *testing.T) {
		inp.SortBy = "REWARD"
		rows := pgxpoolmock.NewRows([]string{"id", "wallet_code", "email", "msisdn", "qty", "transaction",
			"cashback", "reward"}).
			AddRow(int64(1), "CODE_A", "someone@email.id", "628118770510", 1, decimal.NewFromInt(25000),
				decimal.NewFromInt(2500), decimal.NewFromInt(1000)).ToPgxRows()
		pool.EXPECT().Query(ctx, cmd+"c.reward  DESC  limit $3 offset $4", inp.SessionRequest.Id, "1", inp.Limit, inp.Start).Return(rows, nil)
		v, ex := persister.SearchByPartner(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})

	t.Run("should success with text search", func(t *testing.T) {
		inp.SortBy = ""
		inp.TextSearch = "someone"
		rows := pgxpoolmock.NewRows([]string{"id", "wallet_code", "email", "msisdn", "qty", "transaction",
			"cashback", "reward"}).
			AddRow(int64(1), "CODE_A", "someone@email.id", "628118770510", 1, decimal.NewFromInt(25000),
				decimal.NewFromInt(2500), decimal.NewFromInt(1000)).ToPgxRows()
		pool.EXPECT().Query(ctx, `select t.id, t.wallet_code, t.email, t.msisdn, 
			t.qty, t.amount as transaction, c.amount as cashback, c.reward 
			from transactions t inner join cashbacks c 
			on t.kezbek_ref_code = c.kezbek_ref_code
			where t.partner_id = $1  AND (UPPER(t.msisdn) like UPPER($2) OR UPPER(t.email) like UPPER($2) OR UPPER(t.kezbek_ref_code) like UPPER($2) ) `+`
			order by `+" t.id   DESC  limit $3 offset $4", inp.SessionRequest.Id, "%"+inp.TextSearch+"%", inp.Limit, inp.Start).Return(rows, nil)
		v, ex := persister.SearchByPartner(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})

	t.Run("should return exception on failed to query", func(t *testing.T) {
		inp.SortBy = ""
		inp.TextSearch = "someone"
		pool.EXPECT().Query(ctx, `select t.id, t.wallet_code, t.email, t.msisdn, 
			t.qty, t.amount as transaction, c.amount as cashback, c.reward 
			from transactions t inner join cashbacks c 
			on t.kezbek_ref_code = c.kezbek_ref_code
			where t.partner_id = $1  AND (UPPER(t.msisdn) like UPPER($2) OR UPPER(t.email) like UPPER($2) OR UPPER(t.kezbek_ref_code) like UPPER($2) ) `+`
			order by `+" t.id   DESC  limit $3 offset $4", inp.SessionRequest.Id, "%"+inp.TextSearch+"%", inp.Limit, inp.Start).Return(nil, fmt.Errorf("something went wrong"))
		v, ex := persister.SearchByPartner(inp)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})
}

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
	cmd := `INSERT INTO transactions 
		(status, partner_id, partner, wallet_code, msisdn, email,
		qty, amount, partner_ref_code, kezbek_ref_code,
		is_deleted, created_by, created_date)
		VALUES ($1, $2, $3, $4, $5, $6, 
		$7, $8, $9, $10,
		FALSE, $11, NOW()) RETURNING ID`
	t.Run("should success", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"ID"}).AddRow(1).ToPgxRows()
		pool.EXPECT().BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}).Return(tx, nil)
		tx.EXPECT().Commit(ctx).Times(1).Return(nil)
		tx.EXPECT().Rollback(ctx).Times(1).Return(nil)
		tx.EXPECT().QueryRow(context.Background(), cmd,
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
		tx.EXPECT().QueryRow(context.Background(), cmd,
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
