package repository

import (
	"context"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type Transaction struct {
	Pool   storage.Pooler
	Logger *zap.Logger
}

type TransactionPersister interface {
	Add(trx model.Transaction) (*int64, *model.TechnicalError)
}

func NewTransaction(t Transaction) TransactionPersister {
	return &t
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
