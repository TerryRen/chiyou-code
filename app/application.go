package app

import (
	"time"

	"chiyou.code/mmc/conf"
	"chiyou.code/mmc/lib"
	log "github.com/sirupsen/logrus"
)

type SqlColumn struct {
	// Table Name
	TableName string
	// Column Name
	ColumnName string
	// Column Comment
	ColumnComment *string
	// Column Type
	ColumnType string
	// Column Ordinal Position
	ColumnOrdinal int
	// Column Default
	ColumnDefault *string
	// Is Nullable (YES, NO)
	Nullable bool
	// Column Key (PRI)
	ColumnKey *string
	// Data Type (int, bigint, varchar, tinyint)
	DataType string
	// Character Maximum Length
	CharMaxLength *int
	// Numeric Precision
	NumberPrecision *int
	// Numeric Scale
	NumberScale *int
}

type SqlTable struct {
	// Table Name
	TableName string
	// Table Comment
	TableComment *string
	// Table Columns
	Columns map[string]SqlColumn
}

type DataBaseMeta interface {
	// Fetch tables (contains columns)
	GetTableMetas(dbConfig conf.DatabaseConfig) (tables map[string]SqlTable, err error)
}

// Run java code gen
func RunJava() {
	config := conf.GetAppConfig()

	// rotate log config
	lib.ConfigRotateLogger(
		config.Log.LogLevel,
		config.Log.LogPath,
		config.Log.LogFileName,
		time.Duration(config.Log.LogFileMaxAgeHours)*time.Hour,
		time.Duration(config.Log.LogFileRotationHours)*time.Hour)

	log.Info("[start]")

	if config.Db.DriverName != "mysql" {
		log.Fatalf("unSupported database driver: %v", config.Db.DriverName)
		return
	}

	var dbm DataBaseMeta = &MySqlMeta{}

	log.Info("fetch table metadata start")
	tables, err := dbm.GetTableMetas(config.Db)
	if err != nil {
		log.Error("fetch table metadata with error: ", err)
		return
	}
	log.Infof("fetch table metadata end, total: %v", len(tables))

	for tabName, table := range tables {
		log.Debug("==================================")
		log.Debugf("table: %v", tabName)
		log.Debug("==================================")

		var colComment, colDefault string
		for colName, column := range table.Columns {
			colComment = "NULL"
			colDefault = "NULL"

			if column.ColumnComment != nil {
				colComment = *column.ColumnComment
			}
			if column.ColumnDefault != nil {
				colDefault = *column.ColumnDefault
			}

			log.Debugf("column: %v | %v | %v | %v | %v",
				colName,
				colComment,
				column.ColumnType,
				column.ColumnOrdinal,
				colDefault,
			)
		}
	}
	// Render
	log.Info("[java] render start")
	err = RenderJava(tables, config.Java)
	if err != nil {
		log.Error("[java] render with error: ", err)
		return
	}
	log.Info("[java] render end")
	log.Info("[end]")
}
