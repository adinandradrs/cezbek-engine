package management

import (
	"database/sql"
	"fmt"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/mock/repository"
	"github.com/adinandradrs/cezbek-engine/mock/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestParameter_CacheWallets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	dao, cacher := repository.NewMockParamPersister(ctrl),
		storage.NewMockCacher(ctrl)
	manager := NewParameter(Parameter{
		Logger: logger,
		Cacher: cacher,
		Dao:    dao,
	})
	t.Run("should success", func(t *testing.T) {
		params := []*model.Parameter{{
			ParamGroup: sql.NullString{String: "Group A", Valid: true},
			ParamName:  sql.NullString{String: "Name A", Valid: true},
			ParamValue: sql.NullString{String: "Value A", Valid: true},
		}}
		dao.EXPECT().FindByParamGroup("WALLET_CODE").Return(params, nil)
		cacher.EXPECT().Hset(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)
		ex := manager.CacheWallets()
		assert.Nil(t, ex)
	})
	t.Run("should return exception on DAO failure ops", func(t *testing.T) {
		dao.EXPECT().FindByParamGroup("WALLET_CODE").Return(nil,
			apps.Exception("something went wrong",
				fmt.Errorf("something went wrong"), zap.Any("", ""), logger))
		ex := manager.CacheWallets()
		assert.NotNil(t, ex)
	})
	t.Run("should return exception on Redis failure ops", func(t *testing.T) {
		params := []*model.Parameter{{
			ParamGroup: sql.NullString{String: "Group A", Valid: true},
			ParamName:  sql.NullString{String: "Name A", Valid: true},
			ParamValue: sql.NullString{String: "Value A", Valid: true},
		}}
		dao.EXPECT().FindByParamGroup("WALLET_CODE").Return(params, nil)
		cacher.EXPECT().Hset(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(apps.Exception("something went wrong",
				fmt.Errorf("something went wrong"), zap.Any("", ""), logger))
		ex := manager.CacheWallets()
		assert.NotNil(t, ex)
	})
}

func TestParameter_CacheEmailTemplates(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	dao, cacher := repository.NewMockParamPersister(ctrl),
		storage.NewMockCacher(ctrl)
	manager := NewParameter(Parameter{
		Logger: logger,
		Cacher: cacher,
		Dao:    dao,
	})
	t.Run("should success", func(t *testing.T) {
		params := []*model.Parameter{{
			ParamGroup: sql.NullString{String: "Group A", Valid: true},
			ParamName:  sql.NullString{String: "Name A", Valid: true},
			ParamValue: sql.NullString{String: "Value A", Valid: true},
		}}
		dao.EXPECT().FindByParamGroup("EMAIL_TEMPLATE").Return(params, nil)
		cacher.EXPECT().Hset(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)
		ex := manager.CacheEmailTemplates()
		assert.Nil(t, ex)
	})
	t.Run("should return exception on DAO failure ops", func(t *testing.T) {
		dao.EXPECT().FindByParamGroup("EMAIL_TEMPLATE").Return(nil,
			apps.Exception("something went wrong",
				fmt.Errorf("something went wrong"), zap.Any("", ""), logger))
		ex := manager.CacheEmailTemplates()
		assert.NotNil(t, ex)
	})
	t.Run("should return exception on Redis failure ops", func(t *testing.T) {
		params := []*model.Parameter{{
			ParamGroup: sql.NullString{String: "Group A", Valid: true},
			ParamName:  sql.NullString{String: "Name A", Valid: true},
			ParamValue: sql.NullString{String: "Value A", Valid: true},
		}}
		dao.EXPECT().FindByParamGroup("EMAIL_TEMPLATE").Return(params, nil)
		cacher.EXPECT().Hset(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(apps.Exception("something went wrong",
				fmt.Errorf("something went wrong"), zap.Any("", ""), logger))
		ex := manager.CacheEmailTemplates()
		assert.NotNil(t, ex)
	})
}

func TestParameter_CacheEmailSubjects(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	dao, cacher := repository.NewMockParamPersister(ctrl),
		storage.NewMockCacher(ctrl)
	manager := NewParameter(Parameter{
		Logger: logger,
		Cacher: cacher,
		Dao:    dao,
	})
	t.Run("should success", func(t *testing.T) {
		params := []*model.Parameter{{
			ParamGroup: sql.NullString{String: "Group A", Valid: true},
			ParamName:  sql.NullString{String: "Name A", Valid: true},
			ParamValue: sql.NullString{String: "Value A", Valid: true},
		}}
		dao.EXPECT().FindByParamGroup("EMAIL_SUBJECT").Return(params, nil)
		cacher.EXPECT().Hset(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)
		ex := manager.CacheEmailSubjects()
		assert.Nil(t, ex)
	})
	t.Run("should return exception on DAO failure ops", func(t *testing.T) {
		dao.EXPECT().FindByParamGroup("EMAIL_SUBJECT").Return(nil,
			apps.Exception("something went wrong",
				fmt.Errorf("something went wrong"), zap.Any("", ""), logger))
		ex := manager.CacheEmailSubjects()
		assert.NotNil(t, ex)
	})
	t.Run("should return exception on Redis failure ops", func(t *testing.T) {
		params := []*model.Parameter{{
			ParamGroup: sql.NullString{String: "Group A", Valid: true},
			ParamName:  sql.NullString{String: "Name A", Valid: true},
			ParamValue: sql.NullString{String: "Value A", Valid: true},
		}}
		dao.EXPECT().FindByParamGroup("EMAIL_SUBJECT").Return(params, nil)
		cacher.EXPECT().Hset(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(apps.Exception("something went wrong",
				fmt.Errorf("something went wrong"), zap.Any("", ""), logger))
		ex := manager.CacheEmailSubjects()
		assert.NotNil(t, ex)
	})
}
