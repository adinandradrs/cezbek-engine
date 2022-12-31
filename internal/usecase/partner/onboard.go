package partner

import (
	"github.com/adinandradrs/cezbek-engine/internal/adaptor"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/repository"
	"github.com/adinandradrs/cezbek-engine/internal/storage"
	"go.uber.org/zap"
)

type Onboard struct {
	Dao         repository.PartnerPersister
	CiamWatcher adaptor.CiamWatcher
	Cacher      storage.Cacher
	Logger      *zap.Logger
}

type OnboardManager interface {
	AuthenticateClient(inp *model.ClientAuthenticationRequest) (*model.ClientAuthenticationResponse, *model.BusinessError)
	AuthenticateOfficer(inp *model.OfficerAuthenticationRequest) (*model.OfficerAuthenticationResponse, *model.BusinessError)
	ValidateAuthOfficer(inp *model.OfficerValidationRequest) (*model.OfficerValidationResponse, *model.BusinessError)
}

func NewOnboard(o Onboard) OnboardManager {
	return &o
}

func (o Onboard) AuthenticateClient(inp *model.ClientAuthenticationRequest) (*model.ClientAuthenticationResponse, *model.BusinessError) {
	//TODO implement me
	panic("implement me")
}

func (o Onboard) AuthenticateOfficer(inp *model.OfficerAuthenticationRequest) (*model.OfficerAuthenticationResponse, *model.BusinessError) {
	//TODO implement me
	panic("implement me")
}

func (o Onboard) ValidateAuthOfficer(inp *model.OfficerValidationRequest) (*model.OfficerValidationResponse, *model.BusinessError) {
	//TODO implement me
	panic("implement me")
}
