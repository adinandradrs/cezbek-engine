package h2h

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"time"
)

func NewLinksaja(linksaja Linksaja) FactoryProvider {
	return &linksaja
}

func (l *Linksaja) SendCashback(inp *model.H2HSendCashbackRequest) (*model.TransactionResponse, *model.BusinessError) {
	token, ex := l.Cacher.Get("H2H:LINKSAJA", "TOKEN")
	if ex != nil {
		auth, ex := l.Authorization()
		if ex != nil {
			return nil, &model.BusinessError{
				ErrorCode:    apps.ErrCodeBussH2HCashbackFailed,
				ErrorMessage: apps.ErrMsgBussH2HCashbackFailed,
			}
		}
		l.Cacher.Set("H2H:LINKSAJA", "TOKEN", auth.Token, l.TokenTTL)
		token = auth.Token
	}
	v, ex := l.FundTransfer(&model.LinksajaFundTransferRequest{
		Bearer: token,
		Msisdn: inp.Destination,
		Amount: inp.Amount,
		Notes:  inp.Notes,
	})
	if ex != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeBussH2HCashbackFailed,
			ErrorMessage: apps.ErrMsgBussH2HCashbackFailed,
		}
	}
	return &model.TransactionResponse{
		TransactionId:        v.TransactionID,
		TransactionTimestamp: time.Now().Unix(),
	}, nil
}
