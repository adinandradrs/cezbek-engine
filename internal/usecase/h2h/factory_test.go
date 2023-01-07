package h2h

import (
	"encoding/json"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/mock/adaptor"
	"github.com/adinandradrs/cezbek-engine/mock/storage"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFactory_SendCashback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	d, _ := time.ParseDuration("2h")
	cacher := storage.NewMockCacher(ctrl)
	jadp, lsadp, gpadp, mtadp, xadp := adaptor.NewMockJosvoAdapter(ctrl),
		adaptor.NewMockLinksajaAdapter(ctrl), adaptor.NewMockGopaidAdapter(ctrl), adaptor.NewMockMiddletransAdapter(ctrl),
		adaptor.NewMockXenitAdapter(ctrl)
	svc := NewFactory(Factory{
		Cacher:      cacher,
		Josvo:       Josvo{jadp},
		Linksaja:    Linksaja{d, cacher, lsadp},
		Gopaid:      Gopaid{gpadp},
		Middletrans: Middletrans{mtadp},
		Xenit:       Xenit{xadp},
	})
	inp := &model.H2HSendCashbackRequest{
		Amount:      decimal.NewFromInt(5000),
		Notes:       "from kezbek to someone",
		HostCode:    "JOSVOH2H",
		Destination: "628123456789",
		KezbekRefNo: "KEZBEK-001",
	}
	t.Run("should execute LSAJA H2H", func(t *testing.T) {
		inp.WalletCode = "LSAJA"
		providers := []model.H2HPricingProjection{
			{
				Code:       "LSAJAH2H",
				Provider:   "LinkSaja H2H",
				WalletCode: "LSAJA",
				Fee:        decimal.NewFromInt(750),
			},
		}
		c, _ := json.Marshal(providers)
		cacher.EXPECT().Hget("PROVIDER_FEE", "LSAJA").Return(string(c), nil)
		cacher.EXPECT().Get("H2H:LINKSAJA", "TOKEN").Return("", &model.TechnicalError{
			Exception: "cache is empty",
			Ticket:    "ERR-002",
			Occurred:  time.Now().Unix(),
		})
		lsadp.EXPECT().Authorization().Return(&model.LinksajaAuthorizationResponse{
			Token:          "token-abc",
			ThirdPartyName: "KEZBEK",
		}, nil)
		lsadp.EXPECT().FundTransfer(gomock.Any()).Return(&model.LinksajaFundTransferResponse{
			TransactionTime: "1232456435",
			TransactionID:   "TRX-001",
		}, nil)
		cacher.EXPECT().Set("H2H:LINKSAJA", "TOKEN", "token-abc", d)
		v, ex := svc.SendCashback(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})

	t.Run("should execute JOSVO H2H", func(t *testing.T) {
		inp.WalletCode = "JOSVO"
		providers := []model.H2HPricingProjection{
			{
				Code:       "JOSVOH2H",
				Provider:   "JOSVO H2H",
				WalletCode: "JOSVO",
				Fee:        decimal.NewFromInt(900),
			},
		}
		c, _ := json.Marshal(providers)
		cacher.EXPECT().Hget("PROVIDER_FEE", "JOSVO").Return(string(c), nil)
		jadp.EXPECT().AccountTransfer(gomock.Any()).Return(&model.JosvoAccountTransferResponse{
			Notes: "something success",
			Code:  "200",
		}, nil)
		v, ex := svc.SendCashback(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})

	t.Run("should execute GoPaid H2H", func(t *testing.T) {
		inp.WalletCode = "GOPAID"
		providers := []model.H2HPricingProjection{
			{
				Code:       "GOPAIDH2H",
				Provider:   "GoPaid H2H",
				WalletCode: "GOPAID",
				Fee:        decimal.NewFromInt(750),
			},
		}
		c, _ := json.Marshal(providers)
		cacher.EXPECT().Hget("PROVIDER_FEE", "GOPAID").Return(string(c), nil)
		gpadp.EXPECT().Topup(gomock.Any()).Return(&model.GopaidTopupResponse{
			RefCode:   "REF-001",
			Timestamp: "1125642689",
		}, nil)
		v, ex := svc.SendCashback(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})

	t.Run("should execute Middletrans H2H", func(t *testing.T) {
		inp.WalletCode = "GOPAID"
		providers := []model.H2HPricingProjection{
			{
				Code:       "MTRANS",
				Provider:   "Middletrans H2H",
				WalletCode: "GOPAID",
				Fee:        decimal.NewFromInt(750),
			},
		}
		c, _ := json.Marshal(providers)
		cacher.EXPECT().Hget("PROVIDER_FEE", "GOPAID").Return(string(c), nil)
		mtadp.EXPECT().WalletTransfer(gomock.Any()).Return(&model.MiddletransWalletTransferResponse{
			StatusCode:     "00",
			TransactionRef: "TRX-001",
			Message:        "Something",
			IsSuccess:      true,
		}, nil)
		v, ex := svc.SendCashback(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})

	t.Run("should execute Xenit H2H", func(t *testing.T) {
		inp.WalletCode = "GOPAID"
		providers := []model.H2HPricingProjection{
			{
				Code:       "XENIT",
				Provider:   "Xenit H2H",
				WalletCode: "GOPAID",
				Fee:        decimal.NewFromInt(750),
			},
		}
		c, _ := json.Marshal(providers)
		cacher.EXPECT().Hget("PROVIDER_FEE", "GOPAID").Return(string(c), nil)
		xadp.EXPECT().WalletTopup(gomock.Any()).Return(&model.XenitWalletTopupResponse{
			TopupRef:     "REF-001",
			TopupTime:    "1125642689",
			TopupMessage: "success",
			TopupStatus:  "200",
		}, nil)
		v, ex := svc.SendCashback(inp)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})

	t.Run("should return invalid wallet", func(t *testing.T) {
		inp.WalletCode = "XPAY"
		cacher.EXPECT().Hget("PROVIDER_FEE", "XPAY").Return("", &model.TechnicalError{
			Exception: "no cache",
			Occurred:  time.Now().Unix(),
			Ticket:    "ERR-003",
		})
		v, ex := svc.SendCashback(inp)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
	})
}
