package app

import (
	"database/sql"
	"sort"

	"chiyou.code/mmc/conf"
	_ "github.com/go-sql-driver/mysql"
)

type MySqlMeta struct {
}

// fetch tables (contains columns)
func (m *MySqlMeta) GetTableMetas(dbConfig conf.DatabaseConfig) (tables map[string]*SqlTable, err error) {
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

	tables = make(map[string]*SqlTable)

	// Loop over the tables
	for tabs.Next() {
		var table SqlTable = SqlTable{DbType: DB_TYPE_MYSQL, Columns: make(map[string]*SqlColumn)}
		err = tabs.Scan(&table.TableName, &table.TableComment)
		if err != nil {
			return nil, err
		}
		tables[table.TableName] = &table
	}

	//  Loop over the columns
	for cols.Next() {
		var column SqlColumn
		var columnNullable string
		err = cols.Scan(&column.TableName, &column.ColumnName, &column.ColumnComment, &column.ColumnType, &column.ColumnOrdinal,
			&column.ColumnDefault, &columnNullable, &column.ColumnKey, &column.DataType, &column.CharMaxLength, &column.NumberPrecision, &column.NumberScale)
		if err != nil {
			return nil, err
		}
		column.Nullable = columnNullable == "YES"

		// Set table column
		if tab, ok := tables[column.TableName]; ok {
			tab.Columns[column.ColumnName] = &column
		}
	}
	cols.Close()

	// Sort table column by ordinal
	for _, table := range tables {
		table.OrdinalColumns = make([]*SqlColumn, len(table.Columns))

		columnNames := make([]string, 0, len(table.Columns))
		for k := range table.Columns {
			columnNames = append(columnNames, k)
		}
		// Sort the tab by column ordinal
		sort.SliceStable(columnNames, func(i, j int) bool {
			return table.Columns[columnNames[i]].ColumnOrdinal < table.Columns[columnNames[j]].ColumnOrdinal
		})
		// Iterate over the sorted tab
		for i, k := range columnNames {
			table.OrdinalColumns[i] = table.Columns[k]
		}
	}

	return tables, nil
}
