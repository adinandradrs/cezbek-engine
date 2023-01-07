package h2h

import (
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/mock/adaptor"
	"github.com/adinandradrs/cezbek-engine/mock/storage"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLinksaja_SendCashback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adapter, cacher := adaptor.NewMockLinksajaAdapter(ctrl), storage.NewMockCacher(ctrl)
	d, _ := time.ParseDuration("10m")
	svc := NewLinksaja(Linksaja{
		LinksajaAdapter: adapter,
		Cacher:          cacher,
		TokenTTL:        d,
	})
	inp := &model.H2HSendCashbackRequest{
		Amount:      decimal.NewFromInt(5000),
		WalletCode:  "LSAJA",
		Notes:       "from kezbek to someone",
		HostCode:    "LSAJAH2H",
		Destination: "628123456789",
		KezbekRefNo: "KEZBEK-001",
	}

	t.Run("should success", func(t *testing.T) {
		adapter.EXPECT().Authorization().Return(&model.LinksajaAuthorizationResponse{
			Token:          "token-abc",
			ThirdPartyName: "KEZBEK",
		}, nil)
		adapter.EXPECT().FundTransfer(&model.LinksajaFundTransferRequest{
			Bearer: "token-abc",
			Msisdn: inp.Destination,
			Amount: inp.Amount,
			Notes:  inp.Notes,
		}).Return(&model.LinksajaFundTransferResponse{
			TransactionTime: "1232456435",
			TransactionID:   "TRX-001",
		}, nil)
		cacher.EXPECT().Get("H2H:LINKSAJA", "TOKEN").Return("", &model.TechnicalError{
			Exception: "cache is empty",
			Ticket:    "ERR-002",
			Occurred:  time.Now().Unix(),
		})
		cacher.EXPECT().Set("H2H:LINKSAJA", "TOKEN", "token-abc", d)
		tx, ex := svc.SendCashback(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, tx)
	})

	t.Run("should return exception on invalid auth", func(t *testing.T) {
		adapter.EXPECT().Authorization().Return(nil, &model.TechnicalError{
			Exception: "something went wrong",
			Ticket:    "ERR-001",
			Occurred:  time.Now().Unix(),
		})
		cacher.EXPECT().Get("H2H:LINKSAJA", "TOKEN").Return("", &model.TechnicalError{
			Exception: "cache is empty",
			Ticket:    "ERR-002",
			Occurred:  time.Now().Unix(),
		})
		tx, ex := svc.SendCashback(inp)
		assert.NotNil(t, ex)
		assert.Nil(t, tx)
	})

	t.Run("should return exception on failed to fund transfer", func(t *testing.T) {
		adapter.EXPECT().Authorization().Return(&model.LinksajaAuthorizationResponse{
			Token:          "token-abc",
			ThirdPartyName: "KEZBEK",
		}, nil)
		adapter.EXPECT().FundTransfer(&model.LinksajaFundTransferRequest{
			Bearer: "token-abc",
			Msisdn: inp.Destination,
			Amount: inp.Amount,
			Notes:  inp.Notes,
		}).Return(nil, &model.TechnicalError{
			Exception: "something went wrong",
			Ticket:    "ERR-001",
			Occurred:  time.Now().Unix(),
		})
		cacher.EXPECT().Get("H2H:LINKSAJA", "TOKEN").Return("", &model.TechnicalError{
			Exception: "cache is empty",
			Ticket:    "ERR-002",
			Occurred:  time.Now().Unix(),
		})
		cacher.EXPECT().Set("H2H:LINKSAJA", "TOKEN", "token-abc", d)
		tx, ex := svc.SendCashback(inp)
		assert.NotNil(t, ex)
		assert.Nil(t, tx)
	})
}
