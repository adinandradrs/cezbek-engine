package management

import (
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/repository"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"go.uber.org/zap"
)

type Parameter struct {
	Dao    repository.ParamPersister
	Cacher storage.Cacher
	Logger *zap.Logger
}

type ParamManager interface {
	CacheWallets() *model.TechnicalError
	CacheEmailTemplates() *model.TechnicalError
	CacheEmailSubjects() *model.TechnicalError
}

func NewParameter(p Parameter) ParamManager {
	return &p
}

func (p Parameter) groupFetchCache(group string) *model.TechnicalError {
	v, ex := p.Dao.FindByParamGroup(group)
	if ex != nil {
		return ex
	}
	for idx := range v {
		ex = p.Cacher.Hset(group, v[idx].ParamName.String, v[idx].ParamValue.String)
		if ex != nil {
			return ex
		}
	}
	return nil
}

func (p Parameter) CacheWallets() *model.TechnicalError {
	return p.groupFetchCache("WALLET_CODE")
}

func (p Parameter) CacheEmailTemplates() *model.TechnicalError {
	return p.groupFetchCache("EMAIL_TEMPLATE")
}

func (p Parameter) CacheEmailSubjects() *model.TechnicalError {
	return p.groupFetchCache("EMAIL_SUBJECT")
}
