package client

import (
	"database/sql"
	"encoding/json"
	"github.com/adinandradrs/cezbek-engine/internal/adaptor"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/repository"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/h2h"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/workflow"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type Transaction struct {
	Dao repository.TransactionPersister
	workflow.TierProvider
	workflow.CashbackProvider
	h2h.Factory
	SqsAdapter                    adaptor.SQSAdapter
	Cacher                        storage.Cacher
	QueueNotificationEmailInvoice *string
	Logger                        *zap.Logger
}

type TransactionProvider interface {
	Add(inp *model.TransactionRequest) (*model.TransactionResponse, *model.BusinessError)
}

func NewTransaction(t Transaction) TransactionProvider {
	return &t
}

func (t *Transaction) Add(inp *model.TransactionRequest) (*model.TransactionResponse, *model.BusinessError) {
	trx := apps.Transaction(inp.Msisdn)
	_, ex := t.Cacher.Hget("WALLET_CODE", inp.MerchantCode)
	if ex != nil {
		t.Logger.Error("failed to add transaction - invalid wallet code", zap.Any("tx", inp))
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeBussMerchantCodeInvalid,
			ErrorMessage: apps.ErrMsgBussMerchantCodeInvalid,
		}
	}
	data := model.Transaction{
		PartnerId:      inp.SessionRequest.Id,
		Partner:        sql.NullString{String: inp.SessionRequest.Fullname, Valid: true},
		WalletCode:     sql.NullString{String: inp.MerchantCode, Valid: true},
		Qty:            inp.Qty,
		Amount:         inp.Amount,
		Msisdn:         sql.NullString{String: inp.Msisdn, Valid: true},
		Email:          sql.NullString{String: inp.Email, Valid: true},
		KezbekRefCode:  sql.NullString{String: trx.TransactionId, Valid: true},
		PartnerRefCode: sql.NullString{String: inp.TransactionReference, Valid: true},
		BaseEntity: model.BaseEntity{
			CreatedBy: sql.NullInt64{Int64: inp.SessionRequest.Id},
		},
	}
	id, ex := t.Dao.Add(data)
	if ex != nil {
		t.Logger.Error("failed to add transaction - data access", zap.Any("tx", inp))
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeBussClientAddTransaction,
			ErrorMessage: apps.ErrMsgBussClientAddTransaction,
		}
	}
	bx := t.processCashback(&data, inp, id)
	if bx != nil {
		return nil, bx
	}
	return &trx, nil
}

func (t *Transaction) processCashback(data *model.Transaction, inp *model.TransactionRequest, id *int64) *model.BusinessError {
	treward, ex := t.TierProvider.Save(&model.TierRequest{
		PartnerId:     data.PartnerId,
		Email:         data.Email.String,
		Msisdn:        data.Msisdn.String,
		TransactionId: *id,
	})
	if ex != nil {
		t.Logger.Error("failed to save tier", zap.Any("tx", inp))
		return &model.BusinessError{
			ErrorCode:    apps.ErrCodeBussRewardFailed,
			ErrorMessage: apps.ErrMsgBussRewardFailed,
		}
	}
	camt, _ := t.CashbackProvider.FindCashbackAmount(&model.FindCashbackRequest{
		Amount: inp.Amount,
		Qty:    inp.Qty,
	})
	t.Logger.Info("", zap.Any("treward", treward), zap.Any("camt", camt))
	pyld := t.sendCashbackRequest(treward,
		camt, *data)
	v, bx := t.Factory.SendCashback(pyld)
	_ = t.queueEmailInvoice(data, *pyld)
	if bx != nil {
		return bx
	}
	t.Logger.Info("", zap.Any("cashback_resp", v))
	return nil
}

func (t *Transaction) sendCashbackRequest(reward *model.WfRewardTierProjection, cashback *model.FindCashbackResponse, d model.Transaction) *model.H2HSendCashbackRequest {
	subTotal := decimal.Zero
	if reward != nil {
		subTotal = subTotal.Add(reward.Reward)
	}
	if cashback != nil {
		subTotal = subTotal.Add(cashback.Amount)
	}
	return &model.H2HSendCashbackRequest{
		Amount: subTotal,
		Notes: "MSISDN : " + d.Msisdn.String +
			" PARTNER : " + d.Partner.String +
			" REF_CODE : " + d.PartnerRefCode.String,
		KezbekRefNo: d.KezbekRefCode.String,
		WalletCode:  d.WalletCode.String,
		Destination: d.Msisdn.String,
	}
}

func (t *Transaction) queueEmailInvoice(tx *model.Transaction, csb model.H2HSendCashbackRequest) *model.BusinessError {
	sbj, _ := t.Cacher.Hget("EMAIL_SUBJECT", "INVOICE")
	tmpl, _ := t.Cacher.Hget("EMAIL_TEMPLATE", "INVOICE")
	tmpl = strings.ReplaceAll(tmpl, "${reference}", tx.KezbekRefCode.String)
	tmpl = strings.ReplaceAll(tmpl, "${msisdn}", tx.Msisdn.String)
	tmpl = strings.ReplaceAll(tmpl, "${email}", tx.Email.String)
	tmpl = strings.ReplaceAll(tmpl, "${walletCode}", tx.WalletCode.String)
	tmpl = strings.ReplaceAll(tmpl, "${partner}", tx.Partner.String)
	tmpl = strings.ReplaceAll(tmpl, "${qty}", strconv.Itoa(tx.Qty))
	tmpl = strings.ReplaceAll(tmpl, "${transactionAmount}", tx.Amount.String())
	tmpl = strings.ReplaceAll(tmpl, "${cashbackAmount}", csb.Amount.String())
	tmpl = strings.ReplaceAll(tmpl, "\n", "")
	tmpl = strings.ReplaceAll(tmpl, "\t", "")
	msg, err := json.Marshal(model.SendEmailRequest{
		Content:     tmpl,
		Subject:     sbj,
		Destination: tx.Email.String,
	})
	if err != nil {
		return &model.BusinessError{
			ErrorCode:    apps.ErrCodeSomethingWrong,
			ErrorMessage: apps.ErrMsgSomethingWrong,
		}
	}
	err = t.SqsAdapter.SendMessage(*t.QueueNotificationEmailInvoice, string(msg))
	if err != nil {
		return &model.BusinessError{
			ErrorCode:    apps.ErrCodeSomethingWrong,
			ErrorMessage: apps.ErrMsgSomethingWrong,
		}
	}
	return nil
}
