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

func TestJosvo_SendCashback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adapter := adaptor.NewMockJosvoAdapter(ctrl)
	svc := NewJosvo(Josvo{
		JosvoAdapter: adapter,
	})
	inp := &model.H2HSendCashbackRequest{
		Amount:      decimal.NewFromInt(5000),
		WalletCode:  "JOSVO",
		Notes:       "from kezbek to someone",
		HostCode:    "JOSVOH2H",
		Destination: "628123456789",
		KezbekRefNo: "KEZBEK-001",
	}
	t.Run("should success", func(t *testing.T) {
		adapter.EXPECT().AccountTransfer(&model.JosvoAccountTransferRequest{
			Amount:        inp.Amount,
			ClientRefCode: inp.KezbekRefNo,
			PhoneNo:       inp.Destination,
		}).Return(&model.JosvoAccountTransferResponse{
			Notes: "something success",
			Code:  "200",
		}, nil)
		tx, ex := svc.SendCashback(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, tx)
	})

	t.Run("should return exception", func(t *testing.T) {
		adapter.EXPECT().AccountTransfer(&model.JosvoAccountTransferRequest{
			Amount:        inp.Amount,
			ClientRefCode: inp.KezbekRefNo,
			PhoneNo:       inp.Destination,
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
