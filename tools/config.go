package tools

import (
	"fmt"
	"github.com/spf13/viper"
)


var Config = viper.New()

func InitConfig() {
	confPath := "config/"

	Config.AddConfigPath(confPath)
	Config.SetConfigName("config")

	if err := Config.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file config: %s\n", err.Error()))
	}

	Log.Warn("config init success.")
}