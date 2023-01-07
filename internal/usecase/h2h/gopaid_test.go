package h2h

import (
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/mock/adaptor"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGopaid_SendCashback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adapter := adaptor.NewMockGopaidAdapter(ctrl)
	svc := NewGopaid(Gopaid{
		GopaidAdapter: adapter,
	})
	inp := &model.H2HSendCashbackRequest{
		Amount:      decimal.NewFromInt(5000),
		WalletCode:  "GOPAID",
		Notes:       "from kezbek to someone",
		HostCode:    "GOPAIDH2H",
		Destination: "628123456789",
		KezbekRefNo: "KEZBEK-001",
	}
	t.Run("should success", func(t *testing.T) {
		adapter.EXPECT().Topup(&model.GopaidTopUpRequest{
			Receipient: inp.Destination,
			AddBalance: inp.Amount,
		}).Return(&model.GopaidTopupResponse{
			RefCode:   "REF-001",
			Timestamp: "1125642689",
		}, nil)
		tx, ex := svc.SendCashback(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, tx)
	})

	t.Run("should return exception", func(t *testing.T) {
		adapter.EXPECT().Topup(&model.GopaidTopUpRequest{
			Receipient: inp.Destination,
			AddBalance: inp.Amount,
		}).Return(nil, &model.TechnicalError{
			Exception: "something went wrong",
			Occurred:  time.Now().Unix(),
			Ticket:    "ERR-001",
		})
		tx, ex := svc.SendCashback(inp)
		assert.NotNil(t, ex)
		assert.Nil(t, tx)
	})
}
