package app

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"chiyou.code/mmc/conf"
	"chiyou.code/mmc/lib/str"
)

const (
	JAVA_TYPE_STRING     = "String"
	JAVA_TYPE_BIGDECIMAL = "BigDecimal"
	JAVA_TYPE_BOOLEAN    = "Boolean"
	JAVA_TYPE_INTEGER    = "Integer"
	JAVA_TYPE_LONG       = "Long"
	JAVA_TYPE_FLOAT      = "Float"
	JAVA_TYPE_DOUBLE     = "Double"
	JAVA_TYPE_BYTE_ARRAY = "byte[]"
	JAVA_TYPE_DATE       = "Date"
	JAVA_TYPE_TIME       = "Time"
	JAVA_TYPE_TIMESTAMP  = "Timestamp"
	JAVA_TYPE_DATE_TIME  = "LocalDateTime"
)

var (
	JAVA_TYPE_MAP = map[string]string{
		"CHAR":      JAVA_TYPE_STRING,
		"VARCHAR":   JAVA_TYPE_STRING,
		"BLOB":      JAVA_TYPE_STRING,
		"TEXT":      JAVA_TYPE_STRING,
		"NUMERIC":   JAVA_TYPE_BIGDECIMAL,
		"DECIMAL":   JAVA_TYPE_BIGDECIMAL,
		"BIT":       JAVA_TYPE_BOOLEAN,
		"BOOL":      JAVA_TYPE_BOOLEAN,
		"BOOLEAN":   JAVA_TYPE_BOOLEAN,
		"TINYINT":   JAVA_TYPE_INTEGER,
		"SMALLINT":  JAVA_TYPE_INTEGER,
		"MEDIUMINT": JAVA_TYPE_INTEGER,
		"INT":       JAVA_TYPE_INTEGER,
		"INTEGER":   JAVA_TYPE_INTEGER,
		"BIGINT":    JAVA_TYPE_LONG,
		"REAL":      JAVA_TYPE_FLOAT,
		"FLOAT":     JAVA_TYPE_FLOAT,
		"DOUBLE":    JAVA_TYPE_DOUBLE,
		"BINARY":    JAVA_TYPE_BYTE_ARRAY,
		"VARBINARY": JAVA_TYPE_BYTE_ARRAY,
		"DATE":      JAVA_TYPE_DATE,
		"TIME":      JAVA_TYPE_TIME,
		"TIMESTAMP": JAVA_TYPE_TIMESTAMP,
		"DATETIME":  JAVA_TYPE_DATE_TIME,
	}
)

// model.tmpl
type Model struct {
	// Base Package
	BasePackage string
	// Table Comment
	TableComment string
	// Table Class Name
	TableClassName string
	// Create Time
	CreateTime string
	// Field Template String
	Fields []string
}

type ModelField struct {
	FieldComment    string
	FieldAnnotation string
	FieldType       string
	FieldName       string
}

// Init template
func initTemplate(javaConf conf.RenderConfig) (t *template.Template, err error) {
	var files []string
	if err = filepath.Walk(javaConf.TemplateFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	t, err = template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Init output folder
func initOutputFolder(javaConf conf.RenderConfig) (err error) {
	// Check if the directory exists
	_, err = os.Stat(javaConf.OutputFolder)
	if errors.Is(err, os.ErrNotExist) {
		// Create the directory if it does not exist
		err = os.MkdirAll(javaConf.OutputFolder, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// Init Class Type (Field Type & Field Name)
func initClassFieldType(javaConf conf.RenderConfig, tab SqlTable) (err error) {
	// Class
	rep := strings.NewReplacer(javaConf.IgnorePrefix, "", javaConf.IgnoreSuffix, "")
	tab.TableClassName = str.UnderscoreToCamel(rep.Replace(tab.TableName))
	// Field
	for _, column := range tab.Columns {

		column.ClassFieldName = ""
		column.ClassFieldType = ""
	}

	return nil
}

// Render model object
func renderModel(t *template.Template, javaConf conf.RenderConfig, tab SqlTable) (err error) {
	var model Model = Model{
		BasePackage:    javaConf.BasePackage,
		TableComment:   tab.TableName,
		TableClassName: tab.TableClassName,
		CreateTime:     time.Now().Format("2006-01-02 15:04:05"),
		Fields:         make([]string, len(tab.Columns)),
	}
	if tab.TableComment != nil {
		model.TableComment = *tab.TableComment
	}

	buf := &bytes.Buffer{}
	for colName, column := range tab.Columns {
		field := ModelField{
			FieldComment:    colName,
			FieldAnnotation: "@Min(value = 0)",
			FieldType:       column.ClassFieldType,
			FieldName:       column.ClassFieldName,
		}
		if column.ColumnComment != nil {
			field.FieldComment = *column.ColumnComment
		}

		err = t.ExecuteTemplate(buf, "model.field.tmpl", field)
		if err != nil {
			return err
		}
		model.Fields = append(model.Fields, buf.String())
		buf.Reset()
	}

	// Create a text file to write the output
	f, err := os.Create(filepath.Join(javaConf.OutputFolder, model.TableClassName+".java"))
	if err != nil {
		return err
	}
	defer f.Close()

	err = t.ExecuteTemplate(f, "model.tmpl", model)
	if err != nil {
		return err
	}
	return nil
}

func RenderJava(tables map[string]SqlTable, javaConf conf.RenderConfig) (err error) {
	t, err := initTemplate(javaConf)
	if err != nil {
		return err
	}
	err = initOutputFolder(javaConf)
	if err != nil {
		return err
	}
	// Loop over the tables
	for _, tab := range tables {
		// TODO Exclude pattern filter
		// Init Class Type
		initClassFieldType(javaConf, tab)
		// Render Model File
		err = renderModel(t, javaConf, tab)
		if err != nil {
			return err
		}
	}
	return nil
}
