package h2h

import (
	"encoding/json"
	"github.com/adinandradrs/cezbek-engine/internal/adaptor"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"strings"
	"time"
)

type (
	Factory struct {
		storage.Cacher
		Linksaja
		Josvo
		Gopaid
		Middletrans
		Xenit
	}

	Linksaja struct {
		TokenTTL time.Duration
		storage.Cacher
		adaptor.LinksajaAdapter
	}

	Josvo struct {
		adaptor.JosvoAdapter
	}

	Gopaid struct {
		adaptor.GopaidAdapter
	}

	Middletrans struct {
		adaptor.MiddletransAdapter
	}

	Xenit struct {
		adaptor.XenitAdapter
	}
)

type FactoryProvider interface {
	SendCashback(inp *model.H2HSendCashbackRequest) (*model.TransactionResponse, *model.BusinessError)
}

func NewFactory(f Factory) Factory {
	return f
}

func (f Factory) SendCashback(inp *model.H2HSendCashbackRequest) (*model.H2HTransactionResponse, *model.BusinessError) {
	var factory FactoryProvider
	w := strings.ToUpper(inp.WalletCode)
	v, ex := f.Cacher.Hget("PROVIDER_FEE", w)
	if ex != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeBussMerchantCodeInvalid,
			ErrorMessage: apps.ErrMsgBussMerchantCodeInvalid,
		}
	}
	var providers []model.H2HPricingProjection
	_ = json.Unmarshal([]byte(v), &providers)
	if providers[0].Code == apps.H2HLinksaja {
		factory = NewLinksaja(f.Linksaja)
	} else if providers[0].Code == apps.H2HJosvo {
		factory = NewJosvo(f.Josvo)
	} else if providers[0].Code == apps.H2HGpaid {
		factory = NewGopaid(f.Gopaid)
	} else if providers[0].Code == apps.H2HMidtrans {
		factory = NewMiddletrans(f.Middletrans)
	} else if providers[0].Code == apps.H2HXenit {
		factory = NewXenit(f.Xenit)
	} else {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeBadPayload,
			ErrorMessage: apps.ErrMsgBadPayload,
		}
	}
	trx, bx := factory.SendCashback(inp)
	if bx != nil {
		return nil, bx
	}
	return &model.H2HTransactionResponse{
		TransactionResponse: *trx,
		HostCode:            providers[0].Code,
	}, nil
}
