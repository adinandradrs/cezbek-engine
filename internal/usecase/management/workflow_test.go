package management

import (
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

func TestWorkflow_CacheRewardTiers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)

	dao, cacher := repository.NewMockWorkflowPersister(ctrl), storage.NewMockCacher(ctrl)
	svc := NewWorkflow(Workflow{
		Dao:    dao,
		Logger: logger,
		Cacher: cacher,
	})
	t.Run("should success", func(t *testing.T) {
		ptier := "SILVER"
		pgrade := 2
		dao.EXPECT().FindRewardTiers().Return([]model.WfRewardTierProjection{
			{
				Reward:    decimal.NewFromInt(1000),
				Recurring: 3,
				Tier:      "GOLD",
				Grade:     3,
				PrevTier: model.WfRewardTierGradeProjection{
					Grade: &pgrade,
					Tier:  &ptier,
				},
				MaxRecurring: 7,
			},
		}, nil)
		cacher.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
		ex := svc.CacheRewardTiers()
		assert.Nil(t, ex)
	})

	t.Run("should return exception on failed to query", func(t *testing.T) {
		dao.EXPECT().FindRewardTiers().Return(nil, &model.TechnicalError{
			Exception: "something went wrong",
			Occurred:  time.Now().Unix(),
			Ticket:    "ERR-001",
		})
		ex := svc.CacheRewardTiers()
		assert.NotNil(t, ex)
	})
}
