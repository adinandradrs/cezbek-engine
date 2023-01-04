package h2h

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"strconv"
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
	ts, _ := strconv.ParseInt(v.TopupTime, 10, 64)
	return &model.TransactionResponse{
		TransactionId:        v.TopupRef,
		TransactionTimestamp: ts,
	}, nil
}
