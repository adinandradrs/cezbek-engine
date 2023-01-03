package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/chrisyxlee/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestH2H_Providers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	pool := pgxpoolmock.NewMockPgxIface(ctrl)
	ctx := context.Background()
	persister := NewH2H(H2H{
		Logger: logger,
		Pool:   pool,
	})
	t.Run("should success", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"id", "code", "provider"}).AddRow(
			int64(1),
			sql.NullString{String: "CODEA", Valid: true},
			sql.NullString{String: "Provider A", Valid: true},
		).ToPgxRows()
		pool.EXPECT().Query(ctx, "SELECT id, code, provider from h2h_providers where status = $1 and is_deleted = false", apps.StatusActive).
			Return(rows, nil)
		v, ex := persister.Providers()
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})

	t.Run("should return exception on query error", func(t *testing.T) {
		pool.EXPECT().Query(ctx, "SELECT id, code, provider from h2h_providers where status = $1 and is_deleted = false", apps.StatusActive).
			Return(nil, fmt.Errorf("something went wrong"))
		v, ex := persister.Providers()
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})

	t.Run("should return exception on map result", func(t *testing.T) {
		rows := pgxpoolmock.NewRows([]string{"id", "code", "provider"}).AddRow(
			1,
			sql.NullString{String: "CODEA", Valid: true},
			sql.NullString{String: "Provider A", Valid: true},
		).ToPgxRows()
		pool.EXPECT().Query(ctx, "SELECT id, code, provider from h2h_providers where status = $1 and is_deleted = false", apps.StatusActive).
			Return(rows, nil)
		v, ex := persister.Providers()
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})
}

func TestH2H_Pricelists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	pool := pgxpoolmock.NewMockPgxIface(ctrl)
	ctx := context.Background()
	persister := NewH2H(H2H{
		Logger: logger,
		Pool:   pool,
	})
	cmd := `select wallet_code, array_to_json(array_agg(DISTINCT jsonb_build_object(
		'wallet_code', wallet_code,
		'code', code,
		'provider', provider,
		'fee', fee
	))) prices from (select fees.wallet_code,
		   providers.code,
		   providers.provider,
		   min(fees.fee) as fee
	from h2h_provider_fees fees
			 inner join h2h_providers providers on fees.h2h_provider_id = providers.id
	where fees.is_deleted = false and fees.status = $1 and providers.is_deleted = false and providers.status = $1
	group by fees.h2h_provider_id, fees.wallet_code, providers.provider, providers.code
	order by fees.wallet_code, fee ASC) as price_list
	group by wallet_code`
	t.Run("should success", func(t *testing.T) {
		fee, _ := decimal.NewFromString("750")
		prices := []model.H2HPricingProjection{
			{
				WalletCode: "WALLETOCODE_A",
				Code:       "PROVIDER_A",
				Provider:   "Something A",
				Fee:        fee,
			},
		}
		rows := pgxpoolmock.NewRows([]string{"wallet_code", "prices"}).AddRow(
			"WALLETCODE_A", prices).ToPgxRows()
		pool.EXPECT().Query(ctx, cmd, apps.StatusActive).Return(rows, nil)
		v, ex := persister.Pricelists()
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})
	t.Run("should return exception on query error", func(t *testing.T) {
		pool.EXPECT().Query(ctx, cmd, apps.StatusActive).Return(nil, fmt.Errorf("something went wrong"))
		v, ex := persister.Pricelists()
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})
	t.Run("should return exception on map result", func(t *testing.T) {
		prices := []model.H2HPricingProjection{
			{
				WalletCode: "WALLETOCODE_A",
				Code:       "PROVIDER_A",
				Provider:   "Something A",
			},
		}
		j, _ := json.Marshal(prices)
		rows := pgxpoolmock.NewRows([]string{"wallet_code", "prices"}).AddRow(
			"WALLETCODE_A", string(j)).ToPgxRows()
		pool.EXPECT().Query(ctx, cmd, apps.StatusActive).Return(rows, nil)
		v, ex := persister.Pricelists()
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})
}
