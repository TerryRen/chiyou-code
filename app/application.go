package app

import (
	"database/sql"
	"time"

	"chiyou.code/mmc/conf"
	"chiyou.code/mmc/lib"
	_ "github.com/go-sql-driver/mysql"
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

	getTable(config.Db)

	return nil
}

func getTable(dbConfig conf.DatabaseConfig) (err error) {
	// Open a database connection
	db, err := sql.Open(dbConfig.DriverName, dbConfig.DataSourceName)
	if err != nil {
		return err
	}
	defer db.Close()

	// Execute the statement with parameter values
	rows, err := db.Query("SHOW TABLES FROM cb_er;")
	if err != nil {
		return err
	}

	defer rows.Close()

	// Loop over the rows and print the results
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			return err
		}
		log.Info(tableName)
	}

	return nil
}
