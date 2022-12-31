package partner

import (
	"github.com/adinandradrs/cezbek-engine/internal/adaptor"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/adinandradrs/cezbek-engine/internal/repository"
	"go.uber.org/zap"
)

type Onboard struct {
	Dao         repository.PartnerPersister
	CiamWatcher adaptor.CiamWatcher
	Logger      *zap.Logger
}

type OnboardManager interface {
	AuthenticateClient(inp *model.ClientAuthenticationRequest) (*model.ClientAuthenticationResponse, *model.BusinessError)
	AuthenticateOfficer(inp *model.OfficerAuthenticationRequest) (*model.OfficerAuthenticationResponse, *model.BusinessError)
	ValidateAuthOfficer(inp *model.OfficerValidationRequest) (*model.OfficerValidationResponse, *model.BusinessError)
}
