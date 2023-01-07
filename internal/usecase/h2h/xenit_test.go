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

func TestXenit_SendCashback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adapter := adaptor.NewMockXenitAdapter(ctrl)
	svc := NewXenit(Xenit{
		XenitAdapter: adapter,
	})
	inp := &model.H2HSendCashbackRequest{
		Amount:      decimal.NewFromInt(5000),
		WalletCode:  "XENIT",
		Notes:       "from kezbek to someone",
		HostCode:    "XENITH2H",
		Destination: "628123456789",
		KezbekRefNo: "KEZBEK-001",
	}
	t.Run("should success", func(t *testing.T) {
		adapter.EXPECT().WalletTopup(&model.XenitWalletTopupRequest{
			Wallet:      inp.WalletCode + "_" + inp.Destination,
			Amount:      inp.Amount,
			Beneficiary: inp.Destination,
			RefCode:     inp.KezbekRefNo,
		}).Return(&model.XenitWalletTopupResponse{
			TopupRef:     "REF-001",
			TopupTime:    "1125642689",
			TopupMessage: "success",
			TopupStatus:  "200",
		}, nil)
		tx, ex := svc.SendCashback(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, tx)
	})

	t.Run("should return exception", func(t *testing.T) {
		adapter.EXPECT().WalletTopup(&model.XenitWalletTopupRequest{
			Wallet:      inp.WalletCode + "_" + inp.Destination,
			Amount:      inp.Amount,
			Beneficiary: inp.Destination,
			RefCode:     inp.KezbekRefNo,
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
