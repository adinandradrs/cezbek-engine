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

func TestMiddletrans_SendCashback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adapter := adaptor.NewMockMiddletransAdapter(ctrl)
	svc := NewMiddletrans(Middletrans{
		MiddletransAdapter: adapter,
	})
	inp := &model.H2HSendCashbackRequest{
		Amount:      decimal.NewFromInt(5000),
		WalletCode:  "mtrans",
		Notes:       "from kezbek to someone",
		HostCode:    "MTRANSH2H",
		Destination: "628123456789",
		KezbekRefNo: "KEZBEK-001",
	}

	t.Run("should success", func(t *testing.T) {
		adapter.EXPECT().WalletTransfer(&model.MiddletransWalletTransferRequest{
			Amount:  inp.Amount,
			Account: inp.Destination,
			Wallet:  inp.WalletCode,
		}).Return(&model.MiddletransWalletTransferResponse{
			StatusCode:     "00",
			TransactionRef: "TRX-001",
			Message:        "Something",
			IsSuccess:      true,
		}, nil)
		tx, ex := svc.SendCashback(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, tx)
	})

	t.Run("should return exception", func(t *testing.T) {
		adapter.EXPECT().WalletTransfer(&model.MiddletransWalletTransferRequest{
			Amount:  inp.Amount,
			Account: inp.Destination,
			Wallet:  inp.WalletCode,
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
