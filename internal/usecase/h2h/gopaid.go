package h2h

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
)

func NewGopaid(gopaid Gopaid) FactoryProvider {
	return &gopaid
}

func (g *Gopaid) SendCashback(inp *model.H2HSendCashbackRequest) (*model.TransactionResponse, *model.BusinessError) {
	v, ex := g.Topup(&model.GopaidTopUpRequest{
		Receipient: inp.Destination,
		AddBalance: inp.Amount,
	})
	if ex != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeBussH2HCashbackFailed,
			ErrorMessage: apps.ErrMsgBussH2HCashbackFailed,
		}
	}
	return &model.TransactionResponse{
		TransactionId:        v.RefCode,
		TransactionTimestamp: v.Timestamp,
	}, nil
}
