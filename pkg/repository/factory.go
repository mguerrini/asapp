package repository

import (
	"github.com/challenge/pkg/modules/config"
	"github.com/challenge/pkg/modules/logger"
	"github.com/challenge/pkg/repository/memory"
	"github.com/challenge/pkg/repository/sql"
	"sync"
)

var repositoryFactoryOnce sync.Once
var repositoryFactoryInstance IRepositoryFactory

type IRepositoryFactory interface {
	CreateUserRepository(sessionName string) IUserRepository
	CreateMessageRepository(sessionName string)  IMessageRepository
}

func RepositoryFactory() IRepositoryFactory {
	repositoryFactoryOnce.Do(func() {
		if repositoryFactoryInstance != nil {
			return
		}

		factoryType, err := config.ConfigurationSingleton().GetString("root.repositories.factory_type")

		if err != nil {
			logger.Error("Error getting repository factory type. Use sql factory.", err)
		}

		if factoryType == "" {
			logger.Info("Repository factory type is not defined. Use sql factory.")
			factoryType = "sql"
		}

		if factoryType == "sql" {
			repositoryFactoryInstance = &sql.SqlRepositoryFactory{}
		} else if factoryType == "memory" {
			repositoryFactoryInstance = &memory.MemoryRepositoryFactory{}
		} else if factoryType == "none" {

		} else {
			panic("Invalid repository factory type")
		}
	})

	return repositoryFactoryInstance
}

func SetRepositoryFactory(factory IRepositoryFactory) {
	repositoryFactoryInstance = factory
}






