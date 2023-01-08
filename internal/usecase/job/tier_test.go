package job

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/mock/repository"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTier_Expire(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := apps.NewLog(false)
	dao := repository.NewMockTierPersister(ctrl)
	exp, _ := time.ParseDuration("1h")
	svc := NewTier(Tier{
		Expired: &exp,
		Dao:     dao,
		Logger:  logger,
	})
	t.Run("should success", func(t *testing.T) {
		cdata := 5
		dao.EXPECT().CountExpire().Return(&cdata, nil)
		dao.EXPECT().Expire(gomock.Any()).Return(nil)
		svc.Expire()
		assert.Equal(t, 5, cdata)
	})

	t.Run("should no execute expire tier", func(t *testing.T) {
		cdata := 0
		dao.EXPECT().CountExpire().Return(&cdata, nil)
		svc.Expire()
		assert.Equal(t, 0, cdata)
	})

}
