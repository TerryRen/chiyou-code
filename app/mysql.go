package app

import (
	"database/sql"

	"chiyou.code/mmc/conf"
	_ "github.com/go-sql-driver/mysql"
)

type MySqlMeta struct {
}

// fetch tables (contains columns)
func (m *MySqlMeta) GetTableMetas(dbConfig conf.DatabaseConfig) (tables map[string]SqlTable, err error) {
	// Open a database connection
	db, err := sql.Open(dbConfig.DriverName, dbConfig.DataSourceName)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Query all tabs
	tabs, err := db.Query(`
SELECT TABLE_NAME, TABLE_COMMENT 
FROM INFORMATION_SCHEMA.TABLES t 
WHERE TABLE_SCHEMA = ?
	AND TABLE_TYPE = ? 
ORDER BY TABLE_NAME ASC
`, dbConfig.DataBase, "BASE TABLE")
	if err != nil {
		return nil, err
	}
	defer tabs.Close()

	// Query all cols
	cols, err := db.Query(`
SELECT TABLE_NAME, COLUMN_NAME, COLUMN_COMMENT, COLUMN_TYPE, ORDINAL_POSITION, COLUMN_DEFAULT, IS_NULLABLE, COLUMN_KEY, DATA_TYPE, CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE
FROM INFORMATION_SCHEMA.COLUMNS c 
WHERE TABLE_SCHEMA = ?
ORDER BY TABLE_NAME, ORDINAL_POSITION ASC;
`, dbConfig.DataBase)
	if err != nil {
		return nil, err
	}
	defer cols.Close()

	tables = make(map[string]SqlTable)

	// Loop over the tables
	for tabs.Next() {
		var tab SqlTable = SqlTable{Columns: make(map[string]SqlColumn)}
		err = tabs.Scan(&tab.TableName, &tab.TableComment)
		if err != nil {
			return nil, err
		}
		tables[tab.TableName] = tab
	}

	//  Loop over the columns
	for cols.Next() {
		var col SqlColumn
		var columnNullable string
		err = cols.Scan(&col.TableName, &col.ColumnName, &col.ColumnComment, &col.ColumnType, &col.ColumnOrdinal,
			&col.ColumnDefault, &columnNullable, &col.ColumnKey, &col.DataType, &col.CharMaxLength, &col.NumberPrecision, &col.NumberScale)
		if err != nil {
			return nil, err
		}
		col.Nullable = columnNullable == "YES"

		// Set table column
		if tab, ok := tables[col.TableName]; ok {
			tab.Columns[col.ColumnName] = col
		}
	}
	cols.Close()

	return tables, nil
}
