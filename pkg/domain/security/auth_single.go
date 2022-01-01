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
	GeneratePassword(ctx context.Context, password string) (string, error)
}


func AuthenticationFactory() IAuthenticationFactory {
	authenticationFactoryOnce.Do(func() {
		if authenticationFactoryInstance != nil {
			return
		}

		factoryType, err := config.ConfigurationSingleton().GetString("root.authentication.factory_type")

		if err != nil {
			logger.Error("Error getting authentication factory type. Use database factory.", err)
		}

		if factoryType == "" {
			logger.Info("Authentication factory type is not defined. Use database factory.")
			factoryType = "database"
		}

		if factoryType == "database" {
			authenticationFactoryInstance = &dbAuthenticationFactory{}
		} else if factoryType == "none" {
			authenticationFactoryInstance = &noneAuthenticationFactory{}
		} else {
			panic("Invalid authentication factory type")
		}
	})

	return authenticationFactoryInstance
}

func SetAuthenticationFactory(factory IAuthenticationFactory) {
	authenticationFactoryInstance = factory
}



