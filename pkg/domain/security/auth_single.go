package security

import (
	"context"
	"github.com/challenge/pkg/models"
	"github.com/challenge/pkg/modules/config"
	"github.com/challenge/pkg/modules/logger"
	"sync"
)


var authenticationFactoryOnce sync.Once
var authenticationFactoryInstance IAuthenticationFactory

type IAuthenticationFactory interface {
	Create(sessionName string) IAuthentication
}

type IAuthentication interface {
	Authenticate(ctx context.Context, cred models.Login) error
}


func AuthenticationFactory() IAuthenticationFactory {
	authenticationFactoryOnce.Do(func() {
		if authenticationFactoryInstance != nil {
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

		if factoryType == "database" {
			authenticationFactoryInstance = &dbAuthenticationFactory{}
		} else if factoryType == "none" {
			authenticationFactoryInstance = &noneAuthenticationFactory{}
		} else {
			panic("Invalid token authentication factory type")
		}
	})

	return authenticationFactoryInstance
}

func SetAuthenticationFactory(factory IAuthenticationFactory) {
	authenticationFactoryInstance = factory
}



