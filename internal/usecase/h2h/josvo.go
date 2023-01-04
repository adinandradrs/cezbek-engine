package h2h

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"time"
)

func NewJosvo(josvo Josvo) FactoryProvider {
	return &josvo
}

func (j *Josvo) SendCashback(inp *model.H2HSendCashbackRequest) (*model.TransactionResponse, *model.BusinessError) {
	_, ex := j.AccountTransfer(&model.JosvoAccountTransferRequest{
		Amount:        inp.Amount,
		ClientRefCode: inp.KezbekRefNo,
		PhoneNo:       inp.Destination,
	})
	if ex != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeBussH2HCashbackFailed,
			ErrorMessage: apps.ErrMsgBussH2HCashbackFailed,
		}
	}
	return &model.TransactionResponse{
		TransactionId:        apps.TransactionId(inp.Destination),
		TransactionTimestamp: time.Now().Unix(),
	}, nil
}
