package config


func CreateConfigurationManager (path string, filename string) *ConfigurationManagerImpl {
	output := NewConfigurationManager()

	output.Load(path, filename)

	return output
}



