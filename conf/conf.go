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

type DatabaseConfig struct {
	// mysql,sqlserver
	DriverName string
	// {user}:{password}@tcp({ip}:{port})/{db}
	DataSourceName string
	// {db}
	DataBase string
}

type RenderConfig struct {
	// Template folder
	TemplateFolder string
	// Output folder
	OutputFolder string
	// Base Package
	BasePackage string
	// Ignore Prefix
	IgnorePrefix string
	// Ignore Suffix
	IgnoreSuffix string
	// Ignore Columns
	IgnoreColumns []string
}

type AppConfig struct {
	Log  LogConfig
	Db   DatabaseConfig
	Java RenderConfig
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
