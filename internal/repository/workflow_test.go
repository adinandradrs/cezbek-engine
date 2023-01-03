package repository

import (
	"context"
	"fmt"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/chrisyxlee/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
	t.Run("should success", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"cashback_percentage"}).AddRow(decimal.NewFromFloat(1.2)).
			ToPgxRows()
		pool.EXPECT().Query(ctx, `select cashback_percentage from wf_cashbacks 
		where min_qty >= $1 AND 
		($2 >= min_transaction and $2 <= max_transaction) AND
		is_deleted = false AND status = $3`, qty, trx, apps.StatusActive).
			Return(rows, nil)
		v, ex := persister.FindCashbackByTransaction(qty, trx)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})

	t.Run("should return exception on query", func(t *testing.T) {
		pool.EXPECT().Query(ctx, `select cashback_percentage from wf_cashbacks 
		where min_qty >= $1 AND 
		($2 >= min_transaction and $2 <= max_transaction) AND
		is_deleted = false AND status = $3`, qty, trx, apps.StatusActive).
			Return(nil, fmt.Errorf("something went wrong"))
		v, ex := persister.FindCashbackByTransaction(qty, trx)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})

	t.Run("should return exception on map result", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"cashback_percentage"}).AddRow(1).
			ToPgxRows()
		pool.EXPECT().Query(ctx, `select cashback_percentage from wf_cashbacks 
		where min_qty >= $1 AND 
		($2 >= min_transaction and $2 <= max_transaction) AND
		is_deleted = false AND status = $3`, qty, trx, apps.StatusActive).
			Return(rows, nil)
		v, ex := persister.FindCashbackByTransaction(qty, trx)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})
}

func TestWorkflow_FindRewardByTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	pool := pgxpoolmock.NewMockPgxIface(ctrl)
	ctx := context.Background()
	persister := NewWorkflow(Workflow{
		Pool:   pool,
		Logger: logger,
	})
	tier, recurring, level := "BRONZE", 3, 1
	t.Run("should success", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"reward"}).AddRow(decimal.New(int64(10000), 10)).
			ToPgxRows()
		pool.EXPECT().Query(ctx, `select reward from wf_rewards
		where tier = $1
		and recurring = $2
		and tier_level = $3
		is_deleted = false AND status = $4`, tier, recurring, level, apps.StatusActive).
			Return(rows, nil)
		v, ex := persister.FindRewardByTransaction(tier, recurring, level)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})

	t.Run("should return exception on query", func(t *testing.T) {
		pool.EXPECT().Query(ctx, `select reward from wf_rewards
		where tier = $1
		and recurring = $2
		and tier_level = $3
		is_deleted = false AND status = $4`, tier, recurring, level, apps.StatusActive).
			Return(nil, fmt.Errorf("something went wrong"))
		v, ex := persister.FindRewardByTransaction(tier, recurring, level)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})

	t.Run("should return exception on query", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"reward"}).AddRow(1).
			ToPgxRows()
		pool.EXPECT().Query(ctx, `select reward from wf_rewards
		where tier = $1
		and recurring = $2
		and tier_level = $3
		is_deleted = false AND status = $4`, tier, recurring, level, apps.StatusActive).
			Return(rows, nil)
		v, ex := persister.FindRewardByTransaction(tier, recurring, level)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})
}
