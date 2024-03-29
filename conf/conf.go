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
	// Template base class folder
	TemplateBaseClassFolder string
	// Template map
	TemplateMap map[string]string
	// Output folder
	OutputFolder string
	// Model Output folder (sub)
	ModelSubFolder string
	// Dao Output folder  (sub)
	DaoSubFolder string
	// Service Output folder (sub)
	ServiceSubFolder string
	// Author
	Author string
	// Version
	Version string
	// Base Package
	BasePackage string
	// Model Sub Package
	ModelSubPackage string
	// Dao Sub Package
	DaoSubPackage string
	// Service Sub Package
	ServiceSubPackage string
	// Ignore Prefix
	IgnorePrefix string
	// Ignore Suffix
	IgnoreSuffix string
	// Include Table Regex
	IncludeTableRegexs []string
	// Exclude Table Regex
	ExcludeTableRegexs []string
	// Base Model (class) Ignore Columns
	BaseModeltIgnoreColumns []string
	// Update Ignore Columns
	UpdateStatementIgnoreColumns []string
	// Delete Ignore Columns
	DeleteStatementIgnoreColumns []string
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
