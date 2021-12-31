package config

import (
	"sync"
)

var (
	singleton ConfigurationManager
	once      sync.Once
)

func ConfigurationSingleton() ConfigurationManager {
	once.Do(func() {
		if singleton == nil {
			//creo con la configuracion por default
			config := CreateConfigurationManager("", "default.yml")
			singleton = config
		}
	})

	return singleton;
}


func SetSingleton(inst ConfigurationManager) {
	if inst != nil {
		singleton = inst;
	}
}

func CreateSingleton(path, file string) {
	inst := CreateConfigurationManager(path, file)
	SetSingleton(inst)
}

func JoinSingleton(path, file string)  {
	impl, ok := ConfigurationSingleton().(*ConfigurationManagerImpl)

	if ok {
		impl.Join(path, file)
	}
}
