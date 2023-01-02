package management

import (
	"encoding/json"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/repository"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"go.uber.org/zap"
)

type H2H struct {
	Dao    repository.H2HPersister
	Cacher storage.Cacher
	Logger *zap.Logger
}

type H2HManager interface {
	CacheProviders() *model.TechnicalError
	CachePricelists() *model.TechnicalError
}

func NewH2H(h2h H2H) H2HManager {
	return &h2h
}

func (h *H2H) CacheProviders() *model.TechnicalError {
	v, ex := h.Dao.Providers()
	if ex != nil {
		return ex
	}
	for i := range v {
		cache, _ := json.Marshal(v[i])
		h.Cacher.Hset("PROVIDER", v[i].Code.String, cache)
	}
	return nil
}

func (h *H2H) CachePricelists() *model.TechnicalError {
	v, ex := h.Dao.Pricelists()
	if ex != nil {
		return ex
	}
	for i := range v {
		cache, _ := json.Marshal(v[i].Prices)
		h.Cacher.Hset("PROVIDER_FEE", v[i].WalletCode, cache)
	}
	return nil
}
