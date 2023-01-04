package client

import (
	"encoding/json"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/usecase/h2h"
	"github.com/adinandradrs/cezbek-engine/mock/adaptor"
	"github.com/adinandradrs/cezbek-engine/mock/repository"
	"github.com/adinandradrs/cezbek-engine/mock/storage"
	"github.com/adinandradrs/cezbek-engine/mock/usecase/workflow"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTransaction_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	dao, tierProvider, cashbackProvider, cacher, sqsAdapter := repository.NewMockTransactionPersister(ctrl),
		workflow.NewMockTierProvider(ctrl), workflow.NewMockCashbackProvider(ctrl),
		storage.NewMockCacher(ctrl), adaptor.NewMockSQSAdapter(ctrl)
	josvoAdapter := adaptor.NewMockJosvoAdapter(ctrl)
	gopaidAdapter := adaptor.NewMockGopaidAdapter(ctrl)
	linksajaAdapter := adaptor.NewMockLinksajaAdapter(ctrl)
	mtransAdapter := adaptor.NewMockMiddletransAdapter(ctrl)
	xenitAdapter := adaptor.NewMockXenitAdapter(ctrl)
	queueNotificationEmailInvoice := "mock-queue"
	svc := NewTransaction(Transaction{
		Logger:                        logger,
		Cacher:                        cacher,
		TierProvider:                  tierProvider,
		CashbackProvider:              cashbackProvider,
		SqsAdapter:                    sqsAdapter,
		QueueNotificationEmailInvoice: &queueNotificationEmailInvoice,
		Dao:                           dao,
		Factory: h2h.Factory{
			Cacher: cacher,
			Josvo: h2h.Josvo{
				JosvoAdapter: josvoAdapter,
			},
			Linksaja: h2h.Linksaja{
				Cacher:          cacher,
				LinksajaAdapter: linksajaAdapter,
			},
			Gopaid: h2h.Gopaid{
				GopaidAdapter: gopaidAdapter,
			},
			Middletrans: h2h.Middletrans{
				MiddletransAdapter: mtransAdapter,
			},
			Xenit: h2h.Xenit{
				XenitAdapter: xenitAdapter,
			},
		},
	})
	inp := model.TransactionRequest{
		MerchantCode:         "WCODE_A",
		Email:                "someone@email.net",
		TransactionReference: "INV/001/002/003",
		Amount:               decimal.New(150000, 10),
		Qty:                  3,
		Msisdn:               "6281123456890",
		SessionRequest: model.SessionRequest{
			Id:       int64(1),
			Email:    "corporate@email.xyz",
			Fullname: "PT. Corporate A",
		},
	}
	t.Run("should success", func(t *testing.T) {
		providers := []model.H2HPricingProjection{
			{
				Code: "LSAJAH2H",
			},
		}
		b, _ := json.Marshal(providers)
		tierProvider.EXPECT().Save(gomock.Any()).Return(&model.WfRewardTierProjection{Reward: decimal.NewFromInt(100)}, nil)
		cashbackProvider.EXPECT().FindCashbackAmount(gomock.Any()).Return(&model.FindCashbackResponse{
			Amount: decimal.NewFromInt(200),
		}, nil)
		linksajaAdapter.EXPECT().FundTransfer(gomock.Any()).Return(&model.LinksajaFundTransferResponse{
			TransactionID:   "trx-001",
			TransactionTime: "123456",
		}, nil)
		cacher.EXPECT().Get("H2H:LINKSAJA", "TOKEN").Return("something-abc", nil)
		cacher.EXPECT().Hget("WALLET_CODE", inp.MerchantCode).Return("WCODE_A", nil)
		cacher.EXPECT().Hget("PROVIDER_FEE", gomock.Any()).Return(string(b), nil)
		cacher.EXPECT().Hget("EMAIL_SUBJECT", "INVOICE").Return("A subject", nil)
		cacher.EXPECT().Hget("EMAIL_TEMPLATE", "INVOICE").Return("The content", nil)
		sqsAdapter.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return(nil)
		tid := int64(1)
		dao.EXPECT().Add(gomock.Any()).Return(&tid, nil)
		v, ex := svc.Add(&inp)
		assert.Nil(t, ex)
		assert.NotNil(t, v)
	})

	t.Run("should error on data access failed to insert", func(t *testing.T) {
		cacher.EXPECT().Hget("WALLET_CODE", inp.MerchantCode).Return("WCODE_A", nil)
		dao.EXPECT().Add(gomock.Any()).Return(nil, &model.TechnicalError{
			Exception: "something went wrong",
			Occurred:  time.Now().Unix(),
			Ticket:    "ERR-001",
		})
		v, ex := svc.Add(&inp)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
		assert.Equal(t, apps.ErrCodeBussClientAddTransaction, ex.ErrorCode)
	})

	t.Run("should return exception on wallet is not found", func(t *testing.T) {
		cacher.EXPECT().Hget("WALLET_CODE", inp.MerchantCode).Return("", &model.TechnicalError{
			Exception: "something went wrong",
			Occurred:  time.Now().Unix(),
			Ticket:    "ERR-001",
		})
		v, ex := svc.Add(&inp)
		assert.NotNil(t, ex)
		assert.Nil(t, v)
		assert.Equal(t, apps.ErrCodeBussMerchantCodeInvalid, ex.ErrorCode)
	})
}
