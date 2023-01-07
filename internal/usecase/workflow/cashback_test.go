package workflow

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/mock/repository"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCashback_FindCashbackAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)

	dao := repository.NewMockWorkflowPersister(ctrl)
	inp := &model.FindCashbackRequest{
		Qty:    1,
		Amount: decimal.NewFromInt(1000),
	}
	svc := NewCashback(Cashback{
		Logger: logger,
		Dao:    dao,
	})
	t.Run("should success", func(t *testing.T) {
		d := decimal.NewFromFloat(1.5)
		dao.EXPECT().FindCashbackByTransaction(inp.Qty, inp.Amount).Return(&d, nil)
		v, ex := svc.FindCashbackAmount(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})

	t.Run("should return error no cashback", func(t *testing.T) {
		dao.EXPECT().FindCashbackByTransaction(inp.Qty, inp.Amount).Return(nil, &model.TechnicalError{
			Exception: "something went wrong",
			Occurred:  time.Now().Unix(),
			Ticket:    "ERR-001",
		})
		v, ex := svc.FindCashbackAmount(inp)
		assert.Equal(t, apps.ErrCodeBussNoCashback, ex.ErrorCode)
		assert.Nil(t, v)
	})
}
