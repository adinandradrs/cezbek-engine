package partner

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/repository"
	"go.uber.org/zap"
)

type Transaction struct {
	Dao    repository.TransactionPersister
	Logger *zap.Logger
}

type TransactionProvider interface {
	Search(inp *model.SearchRequest) (*model.PartnerTransactionSearchResponse, *model.BusinessError)
	Detail(inp *model.FindByIdRequest) (*model.PartnerTransactionProjection, *model.BusinessError)
}

func NewTransaction(t Transaction) TransactionProvider {
	return &t
}

func (t *Transaction) Search(inp *model.SearchRequest) (*model.PartnerTransactionSearchResponse, *model.BusinessError) {
	model.Page(inp)
	c, countEx := t.Dao.CountByPartner(inp)
	v, searchEx := t.Dao.SearchByPartner(inp)
	if countEx != nil || searchEx != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeNotFound,
			ErrorMessage: apps.ErrMsgNotFound,
		}
	}
	return &model.PartnerTransactionSearchResponse{
		Transactions:       v,
		PaginationResponse: model.Pagination(*c, inp.Limit, inp.Start),
	}, nil
}

func (t *Transaction) Detail(inp *model.FindByIdRequest) (*model.PartnerTransactionProjection, *model.BusinessError) {
	v, ex := t.Dao.DetailByPartner(inp)
	if ex != nil {
		return nil, &model.BusinessError{
			ErrorCode:    apps.ErrCodeNotFound,
			ErrorMessage: apps.ErrMsgNotFound,
		}
	}
	return v, nil
}
