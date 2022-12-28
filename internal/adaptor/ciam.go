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
	"time"
)

type CiamWatcher interface {
	JwtInfo(t string) (map[string]interface{}, *model.TechnicalError)
	OnboardPartner(m model.CiamSignUpPartnerRequest) (*model.CiamUserResponse, *model.TechnicalError)
	Authenticate(m model.CiamSignInRequest) (*model.CiamUserResponse, *model.TechnicalError)
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

func pubKey(token *jwt.Token, url string) (pkey interface{}, err error) {
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

func (c Cognito) parseToken(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, result := token.Method.(*jwt.SigningMethodRSA); !result {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return pubKey(token, c.JWK)
	})
}

func (c Cognito) secretHash(u string) string {
	mac := hmac.New(sha256.New, []byte(c.Scrt))
	mac.Write([]byte(u + c.ClientId))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (c Cognito) JwtInfo(t string) (map[string]interface{}, *model.TechnicalError) {
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

func (c Cognito) SignUpPartner(m model.CiamSignUpPartnerRequest) (*model.CiamUserResponse, *model.TechnicalError) {
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
				Name:  aws.String("phone_number"),
				Value: aws.String(m.PhoneNumber),
			},
		},
	}
	o, err := c.Provider.SignUp(inp)
	if err != nil {
		return nil, apps.Exception("failed to sign up partner", err, zap.String("username", m.Username), c.Logger)
	}
	c.Logger.Info("CIAM output on sign up partner", zap.Any("output", o))
	return &model.CiamUserResponse{
		TransactionResponse: model.TransactionResponse{
			TransactionId:        apps.TransactionId(m.PhoneNumber),
			TransactionTimestamp: time.Now().Unix(),
		},
		SubId: *o.UserSub,
	}, nil
}
