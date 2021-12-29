package config

import (
	"encoding/json"
	"fmt"
	"github.com/challenge/pkg/models/errors"
	"github.com/olebedev/config"
	"os"
	folder "path"
)

type ConfigurationManager interface {
	Exist(path string) bool
	IsNil(path string) bool
	Clean()
	GetObject(path string, emptyObj interface{}) error
	GetString(path string) (string, error)
	GetInt(path string) (int, error)
	GetBool(path string) (bool, error)

	Set(path string, value interface{}) error
}

type ConfigurationManagerImpl struct {
	cfg *config.Config
}

type configEnvelope struct {
	Root interface{} `json:"Root"`
}

func NewConfigurationManager() *ConfigurationManagerImpl {
	output := &ConfigurationManagerImpl{}

	output.Clean()
	return output
}

func (this *ConfigurationManagerImpl) GetCurrentPath() string {
	path, err := os.Getwd()

	if err != nil {
		fmt.Fprintf(os.Stdout, "Error getting source folder: " + err.Error() + "\n")
		return ""
	}

	return path
}

//Load a new configuration file.
func (this *ConfigurationManagerImpl) Load(path string, file string) {
	cfg := this.doParseFile(path, file)
	this.cfg = cfg
}


//Load a file and join it with the existent configuration
func (this *ConfigurationManagerImpl) Join(path string, file string) {
	toJoinConf := this.doParseFile(path, file)

	if toJoinConf == nil {
		panic("Can not join configuration with file " + file)
	}

	nCfg, err := this.cfg.Extend(toJoinConf)

	if err != nil {
		panic(err.Error())
	}

	this.cfg = nCfg
	fmt.Fprintf(os.Stdout, "Configuration file (%s) loaded\n", file)
}

func (this *ConfigurationManagerImpl) existConfigurationFile(path string, file string) bool {
	f, err := this.validateFile(path, file)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Configuration file (%s) doesn't exists - Error: %v", file, err)
		return false
	}

	return len(f) > 0
}

func (this *ConfigurationManagerImpl) doParseFile(path string, file string) *config.Config {

	var err error
	file, err = this.validateFile(path, file)

	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(os.Stdout, "Reading configuration from: %s\n", file)

	cfg, err := config.ParseYamlFile(file)
	if err != nil {
		panic(err.Error())
	}

	return cfg
}

func (this *ConfigurationManagerImpl) validateFile(path string, file string) (string, error) {
	currPath := this.GetCurrentPath()
	exePath, _ := os.Executable()

	if len(path) == 0 {
		fullPath2, exist2 := this.getConfigurationFilePath(exePath, file)
		if exist2 {
			return fullPath2, nil
		}

		fullPath1, exist1 := this.getConfigurationFilePath(currPath, file)
		if exist1 {
			return fullPath1, nil
		}

		return "", errors.NewInternalServerErrorMsg("Configuration file not exists on paths: '" + currPath + ", " + exePath + "' and sub folder 'configs'" )
	}

	//uso el path
	fullPath3, exist3 := this.getConfigurationFilePath(path, file)

	if exist3 {
		return fullPath3, nil
	}


	p1 := folder.Join(exePath, path)
	p2 := folder.Join(currPath, path)
	p3 := folder.Join(path, exePath)
	p4 := folder.Join(path, currPath)

	fullPath4, exist4 := this.getConfigurationFilePath(p1, file)
	if exist4 {
		return fullPath4, nil
	}
	fullPath5, exist5 := this.getConfigurationFilePath(p2, file)
	if exist5 {
		return fullPath5, nil
	}
	fullPath6, exist6 := this.getConfigurationFilePath(p3, file)
	if exist6 {
		return fullPath6, nil
	}
	fullPath7, exist7 := this.getConfigurationFilePath(p4, file)
	if exist7 {
		return fullPath7, nil
	}

	return "", errors.NewInternalServerErrorMsg("Configuration file not exists on paths: '" + path + ", " + p1 + ", " + p2 + ", " + p3 + ", " + p4 + "' and sub folder 'configs'" )
}

func (this *ConfigurationManagerImpl) getConfigurationFilePath(path, file string) (string, bool) {
	fileAux0 := folder.Join(path, file)
	if this.fileExists(fileAux0) {
		return fileAux0, true
	}

	fileAux1 := folder.Join(path, "..\\configs", file)
	if this.fileExists(fileAux1) {
		return fileAux1, true
	}

	fileAux2 := folder.Join(path, "configs", file)
	if this.fileExists(fileAux2) {
		return fileAux2, true
	}

	return "", false
}


func (this *ConfigurationManagerImpl) fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if info == nil {
		return false
	}
	return !info.IsDir()
}

func (this *ConfigurationManagerImpl) Clean() {
	this.cfg = nil
	emptyCfg, _ := config.ParseJson("{}")

	this.cfg = emptyCfg
}

func (this *ConfigurationManagerImpl) Exist(path string) bool {
	_, err := this.cfg.Get(path)

	return err == nil
}

func (this *ConfigurationManagerImpl) IsNil(path string) bool {
	val, err := this.cfg.Get(path)

	if err != nil {
		return true
	}

	return val.Root == nil
}


// Gets

func (this *ConfigurationManagerImpl) GetInt(path string) (int, error) {
	return this.cfg.Int(path)
}

func (this *ConfigurationManagerImpl) GetString(path string) (string, error) {
	return this.cfg.String(path)
}

func (this *ConfigurationManagerImpl) GetBool(path string) (bool, error) {
	return this.cfg.Bool(path)
}

func (this *ConfigurationManagerImpl) GetObject(path string, configType interface{}) error {
	newConfig, err := this.cfg.Get(path)

	if err != nil {
		return err
	}

	jsonObj, err := config.RenderJson(newConfig)

	if err != nil {
		return err
	}

	//json to configs
	objBytes := []byte(jsonObj)

	env := configEnvelope{Root: configType}
	err = json.Unmarshal(objBytes, &env)

	if err != nil {
		return err
	}

	configType = env.Root

	return nil
}

func (this *ConfigurationManagerImpl) Set(path string, value interface{}) error {
	return this.cfg.Set(path, value)
}
