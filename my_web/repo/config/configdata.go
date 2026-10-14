package config

func GetConfigInter() *configReal {
	return &configReal{}
}

// 对外接口
type ConfigInter interface {
	InitConfig(x interface{}, filepath string) error
	GetConfig()
	// MoniterConfigFile(updateConfigChan chan<- *Config)
}

// 具体实现
type configReal struct {
	svrConfig interface{}
	filepath  string
}
