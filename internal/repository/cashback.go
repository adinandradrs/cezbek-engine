package repository

import (
	"context"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type Cashback struct {
	Pool   storage.Pooler
	Logger *zap.Logger
}

type CashbackPersister interface {
	Add(cashback model.Cashback) *model.TechnicalError
}

func NewCashback(c Cashback) CashbackPersister {
	return &c
}

func (c *Cashback) Add(cashback model.Cashback) *model.TechnicalError {
	tx, err := c.Pool.BeginTx(context.Background(),
		pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return apps.Exception("failed to begin add cashback tx", err, zap.Any("", cashback), c.Logger)
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), `INSERT INTO cashbacks 
		(kezbek_ref_code, amount, reward, wallet_code,
		h2h_code, status, is_deleted, created_by, created_date)
		VALUES ($1, $2, $3, $4, $5, $6, FALSE, $7, NOW())`,
		cashback.KezbekRefCode.String, cashback.Amount.Decimal, cashback.Reward.Decimal,
		cashback.WalletCode.String, cashback.H2HCode.String, apps.StatusInactive,
		cashback.CreatedBy.Int64)
	if err != nil {
		return apps.Exception("failed to add cashback tx", err, zap.Any("", cashback), c.Logger)
	}
	if err = tx.Commit(context.Background()); err != nil {
		c.Logger.Panic("failed to commit add cashback trx", zap.Any("cashback", cashback))
	}
	return nil
}
