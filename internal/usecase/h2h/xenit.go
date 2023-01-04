package h2h

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
)

func NewXenit(xenit Xenit) FactoryProvider {
	return &xenit
}

func (x *Xenit) SendCashback(inp *model.H2HSendCashbackRequest) (*model.TransactionResponse, *model.BusinessError) {
	inp.WalletCode = inp.WalletCode + "_" + inp.Destination
	v, ex := x.WalletTopup(&model.XenitWalletTopupRequest{
		Wallet:      inp.WalletCode,
		Amount:      inp.Amount,
		Beneficiary: inp.Destination,
		RefCode:     inp.KezbekRefNo,
	})
	if ex != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeBussH2HCashbackFailed,
			ErrorMessage: apps.ErrMsgBussH2HCashbackFailed,
		}
	}
	return &model.TransactionResponse{
		TransactionId:        v.TopupRef,
		TransactionTimestamp: v.TopupTime,
	}, nil
}
