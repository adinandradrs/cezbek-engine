package repository

import (
	"context"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type Workflow struct {
	Pool   storage.Pooler
	Logger *zap.Logger
}

type WorkflowPersister interface {
	FindCashbackByTransaction(qty int, trx decimal.Decimal) (*decimal.Decimal, *model.TechnicalError)
	FindRewardByTransaction(t string, r int, l int) (*decimal.Decimal, *model.TechnicalError)
}

func NewWorkflow(w Workflow) WorkflowPersister {
	return &w
}

func (w *Workflow) FindCashbackByTransaction(qty int, trx decimal.Decimal) (*decimal.Decimal, *model.TechnicalError) {
	d := decimal.Zero
	query, err := w.Pool.Query(context.Background(), `select cashback_percentage from wf_cashbacks 
		where min_qty >= $1 AND 
		($2 >= min_transaction and $2 <= max_transaction) AND
		is_deleted = false AND status = $3`, qty, trx, apps.StatusActive)
	if err != nil {
		w.Logger.Info("", zap.Int("qty", qty), zap.Any("trx", trx))
		return nil, apps.Exception("failed to find cashback by transaction", err, zap.Error(err), w.Logger)
	}
	err = pgxscan.ScanOne(&d, query)
	if err != nil {
		w.Logger.Info("", zap.Int("qty", qty), zap.Any("trx", trx))
		return nil, apps.Exception("failed to map find cashback by transaction", err, zap.Error(err), w.Logger)
	}
	return &d, nil
}

func (w *Workflow) FindRewardByTransaction(t string, r int, l int) (*decimal.Decimal, *model.TechnicalError) {
	d := decimal.Zero
	query, err := w.Pool.Query(context.Background(), `select reward from wf_rewards
		where tier = $1
		and recurring = $2
		and tier_level = $3
		is_deleted = false AND status = $4`, t, r, l, apps.StatusActive)
	if err != nil {
		w.Logger.Info("", zap.String("tier", t), zap.Any("recurring", r), zap.Int("level", l))
		return nil, apps.Exception("failed to find reward by transaction", err, zap.Error(err), w.Logger)
	}
	err = pgxscan.ScanOne(&d, query)
	if err != nil {
		w.Logger.Info("", zap.String("tier", t), zap.Any("recurring", r), zap.Int("level", l))
		return nil, apps.Exception("failed to map find reward by transaction", err, zap.Error(err), w.Logger)
	}
	return &d, nil
}
