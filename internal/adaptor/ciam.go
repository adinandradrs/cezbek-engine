package adaptor

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/aws/aws-sdk-go/aws"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/golang-jwt/jwt"
	_ "github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
	"go.uber.org/zap"
)

type CiamWatcher interface {
	JwtInfo(t string) (map[string]interface{}, *model.TechnicalError)
	OnboardPartner(m model.CiamOnboardPartnerRequest) (*model.CiamUserResponse, *model.TechnicalError)
	Authenticate(m model.CiamAuthenticationRequest) (*model.CiamAuthenticationResponse, *model.TechnicalError)
}

type (
	Cognito struct {
		Provider *cognito.CognitoIdentityProvider
		UserPool string
		ClientId string
		Scrt     string
		Region   string
		JWK      string
		Logger   *zap.Logger
	}
)

func NewCognito(c Cognito) CiamWatcher {
	return &c
}

func (c *Cognito) pubKey(token *jwt.Token, url string) (pkey interface{}, err error) {
	kid, exists := token.Header["kid"].(string)
	if !exists {
		return nil, errors.New("kid header does not exists")
	}
	kset, err := jwk.ParseString(url)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve JWK")
	}
	keys, _ := kset.LookupKeyID(kid)
	err = keys.Raw(&pkey)
	if err != nil {
		return nil, fmt.Errorf("parsing error")
	}
	return pkey, nil
}

func (c *Cognito) parseToken(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, result := token.Method.(*jwt.SigningMethodRSA); !result {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return c.pubKey(token, c.JWK)
	})
}

func (c *Cognito) secretHash(u string) string {
	mac := hmac.New(sha256.New, []byte(c.Scrt))
	mac.Write([]byte(u + c.ClientId))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (c *Cognito) JwtInfo(t string) (map[string]interface{}, *model.TechnicalError) {
	token, err := c.parseToken(t)
	if err != nil {
		return nil, apps.Exception("failed to get JWT Info from token", err, zap.String("token", t), c.Logger)
	}
	claims, result := token.Claims.(jwt.MapClaims)
	if result && token.Valid {
		return claims, nil
	}
	return nil, apps.Exception("bad JWT claim", fmt.Errorf("failed to claim JWT"), zap.String("token", t), c.Logger)
}

func (c *Cognito) OnboardPartner(m model.CiamOnboardPartnerRequest) (*model.CiamUserResponse, *model.TechnicalError) {
	inp := &cognito.SignUpInput{
		Username:   aws.String(m.Username),
		Password:   aws.String(m.Password),
		ClientId:   aws.String(c.ClientId),
		SecretHash: aws.String(c.secretHash(m.Username)),
		UserAttributes: []*cognito.AttributeType{
			{
				Name:  aws.String("name"),
				Value: aws.String(m.Name),
			},
			{
				Name:  aws.String("email"),
				Value: aws.String(m.Email),
			},
			{
				Name:  aws.String("picture"),
				Value: aws.String(m.Picture),
			},
			{
				Name:  aws.String("phone_number"),
				Value: aws.String("+" + m.PhoneNumber),
			},
		},
	}
	out, err := c.Provider.SignUp(inp)
	if err != nil {
		return nil, apps.Exception("failed to onboard partner", err, zap.String("username", m.Username), c.Logger)
	}
	c.Logger.Info("CIAM output on onboard partner", zap.Any("output", out))
	return &model.CiamUserResponse{
		TransactionResponse: apps.Transaction(m.PhoneNumber),
		SubId:               *out.UserSub,
	}, nil
}

func (c *Cognito) Authenticate(m model.CiamAuthenticationRequest) (*model.CiamAuthenticationResponse, *model.TechnicalError) {
	inp := &cognito.InitiateAuthInput{
		AuthFlow: aws.String(cognito.AuthFlowTypeUserPasswordAuth),
		AuthParameters: map[string]*string{
			"USERNAME":    aws.String(m.Username),
			"PASSWORD":    aws.String(m.Secret),
			"SECRET_HASH": aws.String(c.secretHash(m.Username)),
		},
		ClientId: aws.String(c.ClientId),
	}
	out, err := c.Provider.InitiateAuth(inp)
	if err != nil {
		return nil, apps.Exception("failed to authenticate user", err, zap.String("username", m.Username), c.Logger)
	}
	return &model.CiamAuthenticationResponse{
		AccessToken:  *out.AuthenticationResult.AccessToken,
		Token:        *out.AuthenticationResult.IdToken,
		RefreshToken: *out.AuthenticationResult.RefreshToken,
		ExpiresIn:    *out.AuthenticationResult.ExpiresIn,
	}, nil
}
