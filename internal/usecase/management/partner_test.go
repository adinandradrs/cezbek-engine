package management

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/mock/adaptor"
	"github.com/adinandradrs/cezbek-engine/mock/repository"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"testing"
	"time"
)

func TestPartner_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	dao, ciamWatcher, s3Watcher, pathS3 :=
		repository.NewMockPartnerPersister(ctrl), adaptor.NewMockCiamWatcher(ctrl),
		adaptor.NewMockS3Watcher(ctrl), "/main"
	svc := NewPartner(Partner{
		Dao:         dao,
		CiamWatcher: ciamWatcher,
		S3Watcher:   s3Watcher,
		PathS3:      &pathS3,
		Logger:      logger,
	})
	inp := model.AddPartnerRequest{
		Partner: "PT. Mock Data",
		Address: "Mock Street on Golang",
		Code:    "MOCK",
		Msisdn:  "628123123456",
		Email:   "mock@email.net",
		Logo: multipart.FileHeader{
			Filename: "something,jpeg",
		},
	}
	t.Run("should success", func(t *testing.T) {
		count := 0
		dao.EXPECT().CountByIdentifier(gomock.Any()).Return(&count, nil)
		s3Watcher.EXPECT().Upload(gomock.Any()).Return(nil, nil)
		ciamWatcher.EXPECT().OnboardPartner(gomock.Any()).Return(&model.CiamUserResponse{
			SubId:               "user-123-456",
			TransactionResponse: apps.Transaction(inp.Msisdn),
		}, nil)
		dao.EXPECT().Add(gomock.Any()).Return(nil)
		trx, ex := svc.Add(&inp)
		assert.NotNil(t, trx)
		assert.Nil(t, ex)
	})

	t.Run("should return exception on identifier exists", func(t *testing.T) {
		count := 1
		dao.EXPECT().CountByIdentifier(gomock.Any()).Return(&count, nil)
		trx, ex := svc.Add(&inp)
		assert.Nil(t, trx)
		assert.NotNil(t, ex)
	})

	t.Run("should return exception on error upload logo", func(t *testing.T) {
		count := 0
		dao.EXPECT().CountByIdentifier(gomock.Any()).Return(&count, nil)
		s3Watcher.EXPECT().Upload(gomock.Any()).Return(nil, &model.TechnicalError{
			Exception: "something went wrong",
			Occurred:  time.Now().Unix(),
			Ticket:    uuid.NewString(),
		})
		trx, ex := svc.Add(&inp)
		assert.Nil(t, trx)
		assert.NotNil(t, ex)
	})

	t.Run("should return exception on error register to CIAM", func(t *testing.T) {
		count := 0
		dao.EXPECT().CountByIdentifier(gomock.Any()).Return(&count, nil)
		s3Watcher.EXPECT().Upload(gomock.Any()).Return(nil, nil)
		ciamWatcher.EXPECT().OnboardPartner(gomock.Any()).Return(nil, &model.TechnicalError{
			Exception: "something went wrong",
			Occurred:  time.Now().Unix(),
			Ticket:    uuid.NewString(),
		})
		trx, ex := svc.Add(&inp)
		assert.Nil(t, trx)
		assert.NotNil(t, ex)
	})

	t.Run("should return exception on error add into database", func(t *testing.T) {
		count := 0
		dao.EXPECT().CountByIdentifier(gomock.Any()).Return(&count, nil)
		s3Watcher.EXPECT().Upload(gomock.Any()).Return(nil, nil)
		ciamWatcher.EXPECT().OnboardPartner(gomock.Any()).Return(&model.CiamUserResponse{
			SubId:               "user-123-456",
			TransactionResponse: apps.Transaction(inp.Msisdn),
		}, nil)
		dao.EXPECT().Add(gomock.Any()).Return(&model.TechnicalError{
			Exception: "something went wrong",
			Occurred:  time.Now().Unix(),
			Ticket:    uuid.NewString(),
		})
		trx, ex := svc.Add(&inp)
		assert.Nil(t, trx)
		assert.NotNil(t, ex)
	})

}
