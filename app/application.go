package app

import (
	"fmt"
	"time"

	"chiyou.code/mmc/conf"
	"chiyou.code/mmc/lib"
	log "github.com/sirupsen/logrus"
)

const (
	DB_TYPE_MYSQL    = "mysql"
	CODE_TYPE_JAVA   = "java"
	CODE_TYPE_CSHARP = "csharp"
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
	// Column Key
	// 1. PRI (the column is a PRIMARY KEY or is one of the columns in a multiple-column PRIMARY KEY)
	// 2. UNI (the column is the first column of a UNIQUE index)
	// 3. MUL (the column is the first column of a nonunique index in which multiple occurrences of a given value are permitted within the column)
	// 4. empty (default column)
	ColumnKey *string
	// Data Type (int, bigint, varchar, tinyint)
	DataType string
	// Character Maximum Length
	CharMaxLength *int
	// Numeric Precision
	NumberPrecision *int
	// Numeric Scale
	NumberScale *int
	// Class Field Name
	ClassFieldName string
	// Class Field Type
	ClassFieldType string
}

type SqlTable struct {
	// Database Type
	DbType string
	// Table Name
	TableName string
	// Table Comment
	TableComment *string
	// Table Class Name
	TableClassName string
	// Table Columns map[name]
	Columns map[string]*SqlColumn
	// Table Columns array
	OrdinalColumns []*SqlColumn
}

type DataBaseMeta interface {
	// Fetch tables (contains columns)
	GetTableMetas(dbConfig conf.DatabaseConfig) (tables map[string]*SqlTable, err error)
}

// Run code gen
func Run(codeType string) {
	config := conf.GetAppConfig()

	// rotate log config
	lib.ConfigRotateLogger(
		config.Log.LogLevel,
		config.Log.LogPath,
		config.Log.LogFileName,
		time.Duration(config.Log.LogFileMaxAgeHours)*time.Hour,
		time.Duration(config.Log.LogFileRotationHours)*time.Hour)

	log.Info("[start]")

	var dbm DataBaseMeta
	switch config.Db.DriverName {
	case DB_TYPE_MYSQL:
		dbm = &MySqlMeta{}
	default:
		log.Fatalf("database driver [%v] unSupported", config.Db.DriverName)
		return
	}

	log.Info("[fetch metadata] start")
	tables, err := dbm.GetTableMetas(config.Db)
	if err != nil {
		log.Error("[fetch metadata] with error: ", err)
		return
	}
	log.Infof("[fetch metadata] end, total: %v", len(tables))

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
	switch codeType {
	case CODE_TYPE_JAVA:
		err = RenderJava(tables, config.Java)
	default:
		err = fmt.Errorf("code type [%v] unSupported", codeType)
	}
	if err != nil {
		log.Error("[render] with error: ", err)
	}
	log.Info("[end]")
}
