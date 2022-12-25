package adaptor

import (
	"errors"
	"fmt"
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/golang-jwt/jwt"
	_ "github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
	"go.uber.org/zap"
)

type CiamWatcher interface {
	JwtInfo(t string) (res map[string]interface{}, e *model.TechnicalError)
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

func NewCognito(c *Cognito) CiamWatcher {
	return c
}

func pubKey(token *jwt.Token, jwkurl string) (pkey interface{}, err error) {
	kid, exists := token.Header["kid"].(string)
	if !exists {
		return nil, errors.New("kid header does not exists")
	}
	kset, err := jwk.ParseString(jwkurl)
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

func (c Cognito) JwtInfo(t string) (res map[string]interface{}, e *model.TechnicalError) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, result := token.Method.(*jwt.SigningMethodRSA); !result {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return pubKey(token, c.JWK)
	})
	if err != nil {
		return nil, apps.Exception("failed to get JWT Info from token", err, zap.String("token", t), c.Logger)
	}
	claims, result := token.Claims.(jwt.MapClaims)
	if result && token.Valid {
		return claims, nil
	}
	return nil, apps.Exception("bad JWT claim", fmt.Errorf("failed to claim JWT"), zap.String("token", t), c.Logger)
}
