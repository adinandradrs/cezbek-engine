package management

import (
	"database/sql"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/mock/repository"
	"github.com/adinandradrs/cezbek-engine/mock/storage"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestH2H_CacheProviders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	dao, cacher := repository.NewMockH2HPersister(ctrl), storage.NewMockCacher(ctrl)
	manager := NewH2H(H2H{
		Logger: logger,
		Cacher: cacher,
		Dao:    dao,
	})
	t.Run("should success", func(t *testing.T) {
		dao.EXPECT().Providers().Return([]model.H2HProvider{
			{
				Provider: sql.NullString{String: "Provider A"},
				Code:     sql.NullString{String: "CODE_A"},
				Id:       int64(1),
			},
		}, nil)
		cacher.EXPECT().Hset(gomock.Any(), gomock.Any(), gomock.Any())
		ex := manager.CacheProviders()
		assert.Nil(t, ex)
	})
	t.Run("should return exception on dao error", func(t *testing.T) {
		dao.EXPECT().Providers().Return(nil, &model.TechnicalError{
			Exception: "something went wrong",
			Occurred:  time.Now().Unix(),
			Ticket:    "ERR-001",
		})
		ex := manager.CacheProviders()
		assert.NotNil(t, ex)
	})
}

func TestH2H_CachePricelists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	dao, cacher := repository.NewMockH2HPersister(ctrl), storage.NewMockCacher(ctrl)
	manager := NewH2H(H2H{
		Logger: logger,
		Cacher: cacher,
		Dao:    dao,
	})
	t.Run("should success", func(t *testing.T) {
		dao.EXPECT().Pricelists().Return([]model.H2HPricingsProjection{
			{
				WalletCode: "WALLETCODE_A",
				Prices: []model.H2HPricingProjection{
					{
						WalletCode: "WALLETCODE_A",
						Code:       "CODE_A",
						Provider:   "Provider A",
						Fee:        decimal.New(int64(1000), 10),
					},
				},
			},
		}, nil)
		cacher.EXPECT().Hset(gomock.Any(), gomock.Any(), gomock.Any())
		ex := manager.CachePricelists()
		assert.Nil(t, ex)
	})
	t.Run("should return exception on dao error", func(t *testing.T) {
		dao.EXPECT().Pricelists().Return(nil, &model.TechnicalError{
			Exception: "something went wrong",
			Occurred:  time.Now().Unix(),
			Ticket:    "ERR-001",
		})
		ex := manager.CachePricelists()
		assert.NotNil(t, ex)
	})
}
