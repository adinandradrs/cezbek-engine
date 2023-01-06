package repository

import (
	"context"
	"fmt"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/chrisyxlee/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWorkflow_FindRewardTiers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	pool := pgxpoolmock.NewMockPgxIface(ctrl)
	ctx := context.Background()
	persister := NewWorkflow(Workflow{
		Pool:   pool,
		Logger: logger,
	})
	cmd := `select grade, tier,reward, recurring,
       (select max(recurring) from wf_rewards r where r.grade = m.grade) max_recurring,
       (select jsonb_build_object('grade',grade,'tier', tier) from wf_rewards p where p.id = m.id - 1) prev_tier,
       (select jsonb_build_object('grade',grade,'tier', tier) from wf_rewards n where n.id = m.id + 1) next_tier
		from wf_rewards m order by grade asc, m.tier_level asc`
	t.Run("should success", func(t *testing.T) {
		ptier := "SILVER"
		pgrade := 3
		ntier := "GOLD"
		ngrade := 2
		rows := pgxpoolmock.NewRows([]string{"grade", "tier", "reward", "recurring",
			"max_recurring", "prev_tier", "next_tier"}).AddRow(1, "GOLD", decimal.NewFromInt(10000),
			2, 3, model.WfRewardTierGradeProjection{Tier: &ptier, Grade: &pgrade},
			model.WfRewardTierGradeProjection{Tier: &ntier, Grade: &ngrade}).
			ToPgxRows()
		pool.EXPECT().Query(ctx, cmd).Return(rows, nil)
		v, ex := persister.FindRewardTiers()
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})

	t.Run("should return exception on failed to query", func(t *testing.T) {
		pool.EXPECT().Query(ctx, cmd).Return(nil, fmt.Errorf("something went wrong"))
		v, ex := persister.FindRewardTiers()
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})

	t.Run("should return exception on failed to map result", func(t *testing.T) {
		ptier := "SILVER"
		pgrade := 3
		ntier := "GOLD"
		ngrade := 2
		rows := pgxpoolmock.NewRows([]string{"grade", "tier", "reward", "recurring",
			"max_recurring", "prev_tier", "next_tier"}).AddRow(1, "GOLD", decimal.NewFromInt(10000),
			" 2 ", " 3 ", model.WfRewardTierGradeProjection{Tier: &ptier, Grade: &pgrade},
			model.WfRewardTierGradeProjection{Tier: &ntier, Grade: &ngrade}).
			ToPgxRows()
		pool.EXPECT().Query(ctx, cmd).Return(rows, nil)
		v, ex := persister.FindRewardTiers()
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})
}

func TestWorkflow_FindCashbackByTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	pool := pgxpoolmock.NewMockPgxIface(ctrl)
	ctx := context.Background()
	persister := NewWorkflow(Workflow{
		Pool:   pool,
		Logger: logger,
	})
	qty, trx := 1, decimal.New(15000, 1)
	cmd := `select cashback_percentage from wf_cashbacks 
		where min_qty >= $1 AND max_qty <= $1 AND
		($2 >= min_transaction and $2 <= max_transaction) AND
		is_deleted = false AND status = $3`
	t.Run("should success", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"cashback_percentage"}).AddRow(decimal.NewFromFloat(1.2)).
			ToPgxRows()
		pool.EXPECT().Query(ctx, cmd, qty, trx, apps.StatusActive).
			Return(rows, nil)
		v, ex := persister.FindCashbackByTransaction(qty, trx)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})

	t.Run("should return exception on query", func(t *testing.T) {
		pool.EXPECT().Query(ctx, cmd, qty, trx, apps.StatusActive).
			Return(nil, fmt.Errorf("something went wrong"))
		v, ex := persister.FindCashbackByTransaction(qty, trx)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})

	t.Run("should return exception on map result", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"cashback_percentage"}).AddRow(1).
			ToPgxRows()
		pool.EXPECT().Query(ctx, cmd, qty, trx, apps.StatusActive).
			Return(rows, nil)
		v, ex := persister.FindCashbackByTransaction(qty, trx)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})
}
