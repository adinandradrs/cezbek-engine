package workflow

import (
	"database/sql"
	"encoding/json"
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

func TestTier_Save(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)

	dao, cacher := repository.NewMockTierPersister(ctrl),
		storage.NewMockCacher(ctrl)
	inp := &model.TierRequest{
		PartnerId:     1,
		Msisdn:        "628118770510",
		Email:         "someone@email.net",
		TransactionId: 1,
	}
	exp, _ := time.ParseDuration("72h")
	svc := NewTier(Tier{
		Dao:            dao,
		Logger:         logger,
		Cacher:         cacher,
		ExpiryDuration: exp,
	})
	t.Run("should success with add ops", func(t *testing.T) {
		dao.EXPECT().FindByPartnerMsisdn(inp.PartnerId, inp.Msisdn).
			Return(nil, &model.TechnicalError{
				Ticket:    "ERR-001",
				Occurred:  time.Now().Unix(),
				Exception: "data is not found",
			})
		dao.EXPECT().Add(gomock.Any()).Return(nil)
		_, ex := svc.Save(inp)
		assert.Nil(t, ex)
	})

	t.Run("should success with next tier on update ops", func(t *testing.T) {
		d := &model.Tier{
			PartnerId:            1,
			TransactionRecurring: 2,
			PrevTier:             sql.NullString{String: "BRONZE"},
			PrevGrade:            1,
			CurrentTier:          sql.NullString{String: "SILVER"},
			CurrentGrade:         2,
			ExpiredDate:          sql.NullTime{Time: time.Now()},
			BaseEntity: model.BaseEntity{
				UpdatedBy: sql.NullInt64{Int64: 1},
			},
		}
		dao.EXPECT().FindByPartnerMsisdn(inp.PartnerId, inp.Msisdn).
			Return(d, nil)
		dao.EXPECT().Update(gomock.Any()).Return(nil)
		ngrade := 3
		ntier := "GOLD"
		cache, _ := json.Marshal(model.WfRewardTierProjection{
			MaxRecurring: 3,
			NextTier: model.WfRewardTierGradeProjection{
				Grade: &ngrade,
				Tier:  &ntier,
			},
			Reward:    decimal.NewFromInt(1000),
			Recurring: 3,
		})
		cacher.EXPECT().Get(gomock.Any(), gomock.Any()).Return(string(cache), nil)
		_, ex := svc.Save(inp)
		assert.Nil(t, ex)
	})
}
