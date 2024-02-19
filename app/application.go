package app

import (
	"time"

	"chiyou.code/mmc/conf"
	"chiyou.code/mmc/lib"
	log "github.com/sirupsen/logrus"
)

func Run() error {
	config := conf.GetAppConfig()

	// rotate log config
	lib.ConfigRotateLogger(
		config.Log.LogLevel,
		config.Log.LogPath,
		config.Log.LogFileName,
		time.Duration(config.Log.LogFileMaxAgeHours)*time.Hour,
		time.Duration(config.Log.LogFileRotationHours)*time.Hour)

	log.Info("[app] start...")

	return nil
}
