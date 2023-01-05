package repository

import (
	"context"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type Transaction struct {
	Pool   storage.Pooler
	Logger *zap.Logger
}

type TransactionPersister interface {
	Add(trx model.Transaction) (*int64, *model.TechnicalError)
	SearchByPartner(inp *model.SearchRequest) ([]model.PartnerTransactionProjection, *model.TechnicalError)
	CountByPartner(inp *model.SearchRequest) (*int, *model.TechnicalError)
	DetailByPartner(inp *model.FindByIdRequest) (*model.PartnerTransactionProjection, *model.TechnicalError)
}

func NewTransaction(t Transaction) TransactionPersister {
	return &t
}

func (t *Transaction) DetailByPartner(inp *model.FindByIdRequest) (*model.PartnerTransactionProjection, *model.TechnicalError) {
	v := model.PartnerTransactionProjection{}
	rows, err := t.Pool.Query(context.Background(), `select t.id, t.wallet_code, t.email, t.msisdn, 
			t.qty, t.amount as transaction, c.amount as cashback, c.reward 
			from transactions t inner join cashbacks c 
			on t.kezbek_ref_code = c.kezbek_ref_code and t.id = $1 and t.partner_id = $2`, inp.Id, inp.SessionRequest.Id)
	if err != nil {
		return nil, apps.Exception("failed to get detail by partner", err, zap.Any("", inp), t.Logger)
	}
	defer rows.Close()

	err = pgxscan.ScanOne(&v, rows)
	if err != nil {
		return nil, apps.Exception("failed to map detail by partner", err, zap.Any("", inp), t.Logger)
	}
	return &v, nil
}

func (t *Transaction) CountByPartner(inp *model.SearchRequest) (*int, *model.TechnicalError) {
	var (
		count int
		row   pgx.Row
	)
	where := " AND '1' = $2 "
	if inp.TextSearch != "" {
		where = ` AND (t.msisdn like $2 OR t.email like $2 OR t.kezbek_ref_code like $2 ) `
	}
	cmd := `select count(t.id) 
			from transactions t inner join cashbacks c 
			on t.kezbek_ref_code = c.kezbek_ref_code
			where t.partner_id = $1 ` + where

	if inp.TextSearch == "" {
		row = t.Pool.QueryRow(context.Background(), cmd, inp.SessionRequest.Id, "1")
	} else {
		row = t.Pool.QueryRow(context.Background(), cmd, inp.SessionRequest.Id, inp.TextSearch)
	}
	err := row.Scan(&count)
	if err != nil {
		return nil, apps.Exception("failed to count partner transaction", err, zap.Any("", inp), t.Logger)
	}
	return &count, nil
}

func (t *Transaction) buildOrder(inp *model.SearchRequest) {
	if inp.SortBy == "" {
		inp.SortBy = " t.id "
	}
	if inp.Sort == "" {
		inp.Sort = " DESC "
	}
	if inp.SortBy == "WALLET_CODE" {
		inp.SortBy = "t.wallet_code"
	}
	if inp.SortBy == "AMOUNT" {
		inp.SortBy = "c.amount"
	}
	if inp.SortBy == "REWARD" {
		inp.SortBy = "c.reward"
	}

}

func (t *Transaction) SearchByPartner(inp *model.SearchRequest) ([]model.PartnerTransactionProjection, *model.TechnicalError) {
	var data []model.PartnerTransactionProjection
	where := " AND '1' = $2 "
	if inp.TextSearch != "" {
		where = ` AND (t.msisdn like $2 OR t.email like $2 OR t.kezbek_ref_code like $2 ) `
	}
	t.buildOrder(inp)
	cmd := `select t.id, t.wallet_code, t.email, t.msisdn, 
			t.qty, t.amount as transaction, c.amount as cashback, c.reward 
			from transactions t inner join cashbacks c 
			on t.kezbek_ref_code = c.kezbek_ref_code
			where t.partner_id = $1 ` + where + `
			order by ` + inp.SortBy + " " + inp.Sort + ` limit $3 offset $4`
	var err error
	if inp.TextSearch != "" {
		err = pgxscan.Select(context.Background(), t.Pool, &data,
			cmd, inp.SessionRequest.Id, inp.TextSearch, inp.Limit, inp.Start)
	} else {
		err = pgxscan.Select(context.Background(), t.Pool, &data,
			cmd, inp.SessionRequest.Id, "1", inp.Limit, inp.Start)
	}
	if err != nil {
		return nil, apps.Exception("failed to search partner transaction", err,
			zap.Any("", inp), t.Logger)
	}
	return data, nil
}

func (t *Transaction) Add(trx model.Transaction) (*int64, *model.TechnicalError) {
	tx, err := t.Pool.BeginTx(context.Background(),
		pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return nil, apps.Exception("failed to begin add kezbek tx", err, zap.Any("", trx), t.Logger)
	}
	defer tx.Rollback(context.Background())
	var tid int64
	err = tx.QueryRow(context.Background(), `INSERT INTO transactions 
		(status, partner_id, partner, wallet_code, msisdn, email,
		qty, amount, partner_ref_code, kezbek_ref_code,
		is_deleted, created_by, created_date)
		VALUES ($1, $2, $3, $4, $5, $6, 
		$7, $8, $9, $10,
		FALSE, $11, NOW()) RETURNING ID`,
		apps.StatusInactive, trx.PartnerId, trx.Partner.String, trx.WalletCode.String, trx.Msisdn.String, trx.Email.String,
		trx.Qty, trx.Amount, trx.PartnerRefCode.String, trx.KezbekRefCode.String,
		trx.CreatedBy.Int64).Scan(&tid)
	if err != nil {
		return nil, apps.Exception("failed to add kezbek tx", err, zap.Any("", trx), t.Logger)
	}
	if err = tx.Commit(context.Background()); err != nil {
		t.Logger.Panic("failed to commit add kezbek trx", zap.Any("tx", trx))
	}
	return &tid, nil
}
