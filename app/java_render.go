package app

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"chiyou.code/mmc/conf"
	"chiyou.code/mmc/lib/str"
	log "github.com/sirupsen/logrus"
)

const (
	// Java Type
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
	// Template
	TEMPLATE_MODEL       = "model.tmpl"
	TEMPLATE_MODEL_FIELD = "model.field.tmpl"
	LF                   = "\n"
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
	FieldComment     string
	FieldAnnotations []string
	FieldType        string
	FieldName        string
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
	// Add user function
	t, err = template.New("java").Funcs(template.FuncMap{
		"sub": func(a, b int) int {
			return a - b
		},
	}).ParseFiles(files...)
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
func initClassFieldType(javaConf conf.RenderConfig, tab *SqlTable) (err error) {
	// Class
	rep := strings.NewReplacer(javaConf.IgnorePrefix, "", javaConf.IgnoreSuffix, "")
	tab.TableClassName = str.UnderscoreToCapitalizeCamel(rep.Replace(tab.TableName))
	// Field
	for _, column := range tab.Columns {
		upperDataType := strings.ToUpper(column.DataType)
		if upperDataType == "TINYINT" && strings.ToUpper(column.ColumnType) == "TINYINT(1)" {
			column.ClassFieldType = JAVA_TYPE_BOOLEAN
		} else {
			if tp, ok := JAVA_TYPE_MAP[upperDataType]; ok {
				column.ClassFieldType = tp
			} else {
				return fmt.Errorf("data type: %v mappping failed", column.DataType)
			}
		}
		column.ClassFieldName = str.UnderscoreToLowerCamel(column.ColumnName)
	}
	return nil
}

// Build annotation property (ex: example = "null")
func buildAnnotationProperty(name string, value string, quote bool) string {
	if quote {
		value = `"` + strings.Trim(value, " ") + `"`
	} else {
		value = strings.Trim(value, " ")
	}
	return fmt.Sprintf("%v = %v", name, value)
}

// Build annotation property (ex: @Schema(description = "package count", example = "null") )
func buildAnnotation(name string, properties []string) string {
	if len(properties) == 0 {
		return name
	}
	return fmt.Sprintf("%v(%v)", name, strings.Join(properties, ", "))
}

func buildStringFieldAnnotation(column *SqlColumn) (annotations []string) {
	// @Schema
	description := column.ColumnName
	if column.ColumnComment != nil {
		description = *column.ColumnComment
	}
	example := ""
	if column.ColumnDefault != nil {
		example = *column.ColumnDefault
	}
	ps1 := []string{
		buildAnnotationProperty("description", description, true),
		buildAnnotationProperty("example", example, true),
	}
	annotations = append(annotations, buildAnnotation("@Schema", ps1))

	// @Size
	ps2 := []string{
		buildAnnotationProperty("min", "0", false),
		buildAnnotationProperty("max", strconv.Itoa(*column.CharMaxLength), false),
	}

	annotations = append(annotations, buildAnnotation("@Size", ps2))
	return
}

func buildDefaultFieldAnnotation(column *SqlColumn) (annotations []string) {
	// @Schema
	description := column.ColumnName
	if column.ColumnComment != nil {
		description = *column.ColumnComment
	}
	example := ""
	if column.ColumnDefault != nil {
		example = *column.ColumnDefault
	}
	ps1 := []string{
		buildAnnotationProperty("description", description, true),
		buildAnnotationProperty("example", example, true),
	}

	annotations = append(annotations, buildAnnotation("@Schema", ps1))
	return
}

// Build Field Annotation
func buildFieldAnnotations(column *SqlColumn) (annotations []string) {
	switch column.ClassFieldType {
	case JAVA_TYPE_STRING:
		annotations = buildStringFieldAnnotation(column)
	default:
		annotations = buildDefaultFieldAnnotation(column)
	}

	// @NotNull
	if !column.Nullable {
		annotations = append(annotations, buildAnnotation("@NotNull", nil))
	}
	return
}

// Render model object
func renderModel(t *template.Template, javaConf conf.RenderConfig, tab *SqlTable) (err error) {
	var model Model = Model{
		BasePackage:    javaConf.BasePackage,
		TableComment:   tab.TableName,
		TableClassName: tab.TableClassName,
		CreateTime:     time.Now().Format("2006-01-02 15:04:05"),
		Fields:         make([]string, 0),
	}
	if tab.TableComment != nil {
		model.TableComment = *tab.TableComment
	}

	buf := &bytes.Buffer{}
	for colName, column := range tab.Columns {
		field := ModelField{
			FieldComment: colName,
			FieldType:    column.ClassFieldType,
			FieldName:    column.ClassFieldName,
		}
		if column.ColumnComment != nil {
			field.FieldComment = *column.ColumnComment
		}
		// Build annotation
		field.FieldAnnotations = buildFieldAnnotations(column)
		// Render java field
		err = t.ExecuteTemplate(buf, TEMPLATE_MODEL_FIELD, field)
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

	err = t.ExecuteTemplate(f, TEMPLATE_MODEL, model)
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
		err = initClassFieldType(javaConf, &tab)
		if err != nil {
			log.Errorf("table [%v] init class field type with error: %v", tab.TableName, err)
			continue
		}
		// Render Model File
		err = renderModel(t, javaConf, &tab)
		if err != nil {
			log.Errorf("table [%v] render model with error: %v", tab.TableName, err)
			continue
		}
	}
	return nil
}
