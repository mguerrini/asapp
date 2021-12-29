package security

import (
	"context"
	"github.com/challenge/pkg/modules/config"
	"github.com/challenge/pkg/modules/logger"
	"sync"
)

type TokenStatus string

const (
	SecurityTokenStatus_OK TokenStatus = "SecurityToken_Ok"
	SecurityTokenStatus_Expired TokenStatus = "SecurityToken_Expired"
	SecurityTokenStatus_Invalid TokenStatus = "SecurityToken_Invalid"
)


var tokenAuthenticationFactoryOnce sync.Once
var tokenAuthenticationFactoryInstance ITokenAuthenticationFactory

type ITokenAuthenticationFactory interface {
	Create() ITokenAuthentication
}

type ITokenAuthentication interface {
	GenerateToken(ctx context.Context, user string) string
	ValidateToken(ctx context.Context, token string) TokenStatus
}

func TokenAuthenticationFactory() ITokenAuthenticationFactory {
	tokenAuthenticationFactoryOnce.Do(func() {
		if tokenAuthenticationFactoryInstance != nil {
			return
		}

		factoryType, err := config.ConfigurationSingleton().GetString("root.token_auth.factory_type")

		if err != nil {
			logger.Error("Error getting token auth factory type. Use jws factory.", err)
		}

		if factoryType == "" {
			logger.Info("Token auth factory type is not defined. Use jws factory.")
			factoryType = "jws"
		}

		if factoryType == "jws" {
			tokenAuthenticationFactoryInstance = &jwsTokenAuthenticationFactory{}
		} else if factoryType == "none" {
			tokenAuthenticationFactoryInstance = &noneTokenAuthenticationFactory{}
		} else {
			panic("Invalid token authentication factory type")
		}
	})

	return tokenAuthenticationFactoryInstance
}

func SetTokenAuthenticationFactory(factory ITokenAuthenticationFactory) {
	tokenAuthenticationFactoryInstance = factory
}



