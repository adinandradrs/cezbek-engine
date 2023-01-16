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
	FindRewardTiers() ([]model.WfRewardTierProjection, *model.TechnicalError)
}

func NewWorkflow(w Workflow) WorkflowPersister {
	return &w
}

func (w *Workflow) FindRewardTiers() ([]model.WfRewardTierProjection, *model.TechnicalError) {
	var d []model.WfRewardTierProjection
	rows, err := w.Pool.Query(context.Background(), `select grade, tier,reward, recurring,
       (select max(recurring) from wf_rewards r where r.grade = m.grade) max_recurring,
       (select jsonb_build_object('grade',grade,'tier', tier) from wf_rewards p where p.id = m.id - 1) prev_tier,
       (select jsonb_build_object('grade',grade,'tier', tier) from wf_rewards n where n.id = m.id + 1) next_tier
		from wf_rewards m order by grade asc, m.tier_level asc`)
	if err != nil {
		return nil, apps.Exception("failed to find rewards tiers", err, zap.Error(err), w.Logger)
	}
	defer rows.Close()

	err = pgxscan.ScanAll(&d, rows)
	if err != nil {
		return nil, apps.Exception("failed to map reward tiers", err, zap.Error(err), w.Logger)
	}
	return d, nil
}

func (w *Workflow) FindCashbackByTransaction(qty int, trx decimal.Decimal) (*decimal.Decimal, *model.TechnicalError) {
	d := decimal.Zero
	rows, err := w.Pool.Query(context.Background(), `select cashback_percentage from wf_cashbacks 
		where $1 >= min_qty AND $1 <= max_qty AND
		($2 >= min_transaction and $2 <= max_transaction) AND
		is_deleted = false AND status = $3`, qty, trx, apps.StatusActive)
	if err != nil {
		w.Logger.Info("", zap.Int("qty", qty), zap.Any("trx", trx))
		return nil, apps.Exception("failed to find cashback by transaction", err, zap.Error(err), w.Logger)
	}
	defer rows.Close()

	err = pgxscan.ScanOne(&d, rows)
	if err != nil {
		w.Logger.Info("", zap.Int("qty", qty), zap.Any("trx", trx))
		return nil, apps.Exception("failed to map find cashback by transaction", err, zap.Error(err), w.Logger)
	}
	return &d, nil
}
