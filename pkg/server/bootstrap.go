package server

import (
	"github.com/challenge/pkg/modules/config"
	"github.com/challenge/pkg/modules/storage"
)

func StartServer() {
	//initialize configuration

	env, err := config.ConfigurationSingleton().GetString("root.startup.environment")
	if err == nil && env != "" {
		config.JoinSingleton("", env + ".yml")
	}
}

func FinishServer() {
	//close connections
	storage.DBManagerSingleton().CloseAll()
}

