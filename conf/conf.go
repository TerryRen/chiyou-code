package conf

import (
	"sync"

	"github.com/spf13/viper"
)

// log config
type LogConfig struct {
	// panic,fatal,error,warning|warn,info,debug,trace
	LogLevel             string
	LogPath              string
	LogFileName          string
	LogFileRotationHours int
	LogFileMaxAgeHours   int
}

type JavaConfig struct {
}

type AppConfig struct {
	Log  LogConfig
	Java JavaConfig
}

// config is a package-level variable that stores the configuration object
var appConfig *AppConfig

// once is a sync.Once variable that ensures the configuration is initialized only once
var once sync.Once

// init is a special function that is called when the package is imported

// GetConfig returns the configuration object
func GetAppConfig() *AppConfig {
	// use sync.Once to initialize the config object only once
	once.Do(func() {
		if err := viper.Unmarshal(&appConfig); err != nil {
			panic(err)
		}
	})
	// return the config object
	return appConfig
}
