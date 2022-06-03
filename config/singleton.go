package config

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"sync"
)

var once sync.Once
var instance config

func GetConfigInstance() *config {

	once.Do(func() {

		path := "$OCH_PATH"
		viper.AddConfigPath(path)
		viper.SetConfigName("app")
		viper.SetConfigType("yaml")
		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}

		viper.Unmarshal(&instance)

		if instance.IsDebug {
			log.Info().
				Str("path", viper.ConfigFileUsed()).
				Interface("dao", instance).
				Msg("Config Init")
		}
	})

	return &instance
}
