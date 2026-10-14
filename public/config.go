package public

import (
	"os"

	"github.com/spf13/viper"
)

func InitConfig() *viper.Viper {
	wkDir, _ := os.Getwd()
	VP = viper.New()
	VP.SetConfigName("config")
	VP.SetConfigType("yml")
	VP.AddConfigPath(wkDir)
	if err := VP.ReadInConfig(); err != nil {
		LOG.Fatal(err)
	}
	return VP
}

func GetConfig(key string) string {
	return VP.GetString(key)
}
