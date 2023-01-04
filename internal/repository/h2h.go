package repository

import (
	"context"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"github.com/georgysavva/scany/pgxscan"
	"go.uber.org/zap"
)

type H2H struct {
	Pool   storage.Pooler
	Logger *zap.Logger
}

type H2HPersister interface {
	Providers() ([]model.H2HProvider, *model.TechnicalError)
	Pricelists() ([]model.H2HPricingsProjection, *model.TechnicalError)
}

func NewH2H(h2h H2H) H2HPersister {
	return &h2h
}

func (h2h *H2H) Providers() ([]model.H2HProvider, *model.TechnicalError) {
	var providers []model.H2HProvider
	rows, err := h2h.Pool.Query(context.Background(),
		"SELECT id, code, provider from h2h_providers where status = $1 and is_deleted = false", apps.StatusActive)

	if err != nil {
		return nil, apps.Exception("failed to fetch providers", err, zap.Error(err), h2h.Logger)
	}
	defer rows.Close()

	err = pgxscan.ScanAll(&providers, rows)
	if err != nil {
		return nil, apps.Exception("failed to map fetch providers", err, zap.Error(err), h2h.Logger)
	}
	return providers, nil
}

func (h2h *H2H) Pricelists() ([]model.H2HPricingsProjection, *model.TechnicalError) {
	var pricings []model.H2HPricingsProjection
	cmd := `select wallet_code, array_to_json(array_agg(DISTINCT jsonb_build_object(
		'wallet_code', wallet_code,
		'code', code,
		'provider', provider,
		'fee', fee
	))) prices from (select fees.wallet_code,
		   providers.code,
		   providers.provider,
		   min(fees.fee) as fee
	from h2h_provider_fees fees
			 inner join h2h_providers providers on fees.h2h_provider_id = providers.id
	where fees.is_deleted = false and fees.status = $1 and providers.is_deleted = false and providers.status = $1
	group by fees.h2h_provider_id, fees.wallet_code, providers.provider, providers.code
	order by fees.wallet_code, fee ASC) as price_list
	group by wallet_code`
	rows, err := h2h.Pool.Query(context.Background(), cmd, apps.StatusActive)
	if err != nil {
		return nil, apps.Exception("failed to fetch price lists", err, zap.Error(err), h2h.Logger)
	}
	defer rows.Close()

	err = pgxscan.ScanAll(&pricings, rows)
	if err != nil {
		return nil, apps.Exception("failed to map price lists", err, zap.Error(err), h2h.Logger)
	}
	return pricings, nil
}
