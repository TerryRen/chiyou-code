package app

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
	// Author
	Author string
	// Version
	Version string
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
func initTemplate(renderConf conf.RenderConfig) (t *template.Template, err error) {
	var files []string
	if err = filepath.Walk(renderConf.TemplateFolder, func(path string, info os.FileInfo, err error) error {
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
func initOutputFolder(renderConf conf.RenderConfig) (err error) {
	// Check if the directory exists
	_, err = os.Stat(renderConf.OutputFolder)
	if errors.Is(err, os.ErrNotExist) {
		// Create the directory if it does not exist
		err = os.MkdirAll(renderConf.OutputFolder, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// Init Class Type (Field Type & Field Name)
func initClassFieldType(renderConf conf.RenderConfig, tab *SqlTable) (err error) {
	// Class
	rep := strings.NewReplacer(renderConf.IgnorePrefix, "", renderConf.IgnoreSuffix, "")
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

// Default @Schema(description="xxx" , example="null")
func buildDefaultSchemaAnnotation(column *SqlColumn) (annotation string) {
	// @Schema
	description := column.ColumnName
	if column.ColumnComment != nil {
		description = *column.ColumnComment
	}
	example := "null"
	if column.ColumnDefault != nil {
		example = *column.ColumnDefault
	}
	ps1 := []string{
		buildAnnotationProperty("description", description, true),
		buildAnnotationProperty("example", example, true),
	}
	annotation = buildAnnotation("@Schema", ps1)

	return
}

// String Annotation
func buildStringFieldAnnotation(column *SqlColumn) (annotations []string) {
	annotations = make([]string, 2)
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
	annotations[0] = buildAnnotation("@Schema", ps1)

	// @Size
	ps2 := []string{
		buildAnnotationProperty("min", "0", false),
		buildAnnotationProperty("max", strconv.Itoa(*column.CharMaxLength), false),
	}
	annotations[1] = buildAnnotation("@Size", ps2)

	return
}

// Integer Annotation
func buildIntegerFieldAnnotation(column *SqlColumn) (annotations []string) {
	annotations = make([]string, 3)
	// @Schema
	annotations[0] = buildDefaultSchemaAnnotation(column)

	// @Min
	ps2 := []string{
		buildAnnotationProperty("value", "0", false),
	}
	annotations[1] = buildAnnotation("@Min", ps2)

	// @Max
	ps3 := []string{
		buildAnnotationProperty("value", "Integer.MAX_VALUE", false),
	}
	annotations[2] = buildAnnotation("@Max", ps3)

	return
}

// Long Annotation
func buildLongFieldAnnotation(column *SqlColumn) (annotations []string) {
	annotations = make([]string, 3)
	// @Schema
	annotations[0] = buildDefaultSchemaAnnotation(column)

	// @Min
	ps2 := []string{
		buildAnnotationProperty("value", "0L", false),
	}
	annotations[1] = buildAnnotation("@Min", ps2)

	// @Max
	ps3 := []string{
		buildAnnotationProperty("value", "Long.MAX_VALUE", false),
	}
	annotations[2] = buildAnnotation("@Max", ps3)

	return
}

// BigDecimal Annotation
func buildBigDecimalFieldAnnotation(column *SqlColumn) (annotations []string) {
	annotations = make([]string, 4)
	min := fmt.Sprintf("%v.%v", "0", str.RepeatString(*column.NumberPrecision, "0"))
	max := fmt.Sprintf("%v.%v",
		str.RepeatString((*column.NumberPrecision)-(*column.NumberScale), "9"),
		str.RepeatString(*column.NumberPrecision, "9"))
	// @Schema
	description := column.ColumnName
	if column.ColumnComment != nil {
		description = *column.ColumnComment
	}
	example := "null"
	if column.ColumnDefault != nil {
		example = *column.ColumnDefault
	}
	ps1 := []string{
		buildAnnotationProperty("description", description, true),
		buildAnnotationProperty("example", example, true),
		buildAnnotationProperty("minimum", min, true),
		buildAnnotationProperty("maximum", max, true),
	}
	annotations[0] = buildAnnotation("@Schema", ps1)

	// @Digits
	ps2 := []string{
		buildAnnotationProperty("integer", strconv.Itoa(*column.NumberPrecision), false),
		buildAnnotationProperty("fraction", strconv.Itoa(*column.NumberScale), false),
	}
	annotations[1] = buildAnnotation("@Digits", ps2)

	// @DecimalMin
	ps3 := []string{
		buildAnnotationProperty("value", min, true),
	}
	annotations[2] = buildAnnotation("@DecimalMin", ps3)

	// @DecimalMax
	ps4 := []string{
		buildAnnotationProperty("value", max, true),
	}
	annotations[3] = buildAnnotation("@DecimalMax", ps4)

	return
}

// Boolean Annotation
func buildBooleanFieldAnnotation(column *SqlColumn) (annotations []string) {
	annotations = make([]string, 1)
	// @Schema
	description := column.ColumnName
	if column.ColumnComment != nil {
		description = *column.ColumnComment
	}
	example := "null"
	if column.ColumnDefault != nil {
		if b, err := strconv.ParseBool(*column.ColumnDefault); err == nil {
			example = strconv.FormatBool(b)
		}
	}
	ps1 := []string{
		buildAnnotationProperty("description", description, true),
		buildAnnotationProperty("example", example, true),
	}
	annotations[0] = buildAnnotation("@Schema", ps1)

	return
}

// Default Annotation
func buildDefaultFieldAnnotation(column *SqlColumn) (annotations []string) {
	annotations = make([]string, 1)
	// @Schema
	annotations[0] = buildDefaultSchemaAnnotation(column)

	return
}

// Build Field Annotations
func buildFieldAnnotations(column *SqlColumn) (annotations []string) {
	switch column.ClassFieldType {
	case JAVA_TYPE_STRING:
		annotations = buildStringFieldAnnotation(column)
	case JAVA_TYPE_INTEGER:
		annotations = buildIntegerFieldAnnotation(column)
	case JAVA_TYPE_LONG:
		annotations = buildLongFieldAnnotation(column)
	case JAVA_TYPE_BIGDECIMAL:
		annotations = buildBigDecimalFieldAnnotation(column)
	case JAVA_TYPE_BOOLEAN:
		annotations = buildBooleanFieldAnnotation(column)
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
func renderModel(t *template.Template, renderConf conf.RenderConfig, tab *SqlTable) (err error) {
	var model Model = Model{
		Author:         renderConf.Author,
		Version:        renderConf.Version,
		BasePackage:    renderConf.BasePackage,
		TableComment:   tab.TableName,
		TableClassName: tab.TableClassName,
		CreateTime:     time.Now().Format("2006-01-02 15:04:05"),
		Fields:         make([]string, 0),
	}
	if tab.TableComment != nil {
		model.TableComment = *tab.TableComment
	}

	columnNames := make([]string, 0, len(tab.Columns))
	for k := range tab.Columns {
		columnNames = append(columnNames, k)
	}
	// Sort the tab by column ordinal
	sort.SliceStable(columnNames, func(i, j int) bool {
		return tab.Columns[columnNames[i]].ColumnOrdinal < tab.Columns[columnNames[j]].ColumnOrdinal
	})
	// Iterate over the sorted tab
	buf := &bytes.Buffer{}
	for _, k := range columnNames {
		column := tab.Columns[k]
		field := ModelField{
			FieldComment: column.ColumnName,
			FieldType:    column.ClassFieldType,
			FieldName:    column.ClassFieldName,
		}
		if column.ColumnComment != nil {
			field.FieldComment = *column.ColumnComment
		}
		// Build annotation
		field.FieldAnnotations = buildFieldAnnotations(column)
		// Render java model.field
		err = t.ExecuteTemplate(buf, TEMPLATE_MODEL_FIELD, field)
		if err != nil {
			return err
		}
		model.Fields = append(model.Fields, buf.String())
		buf.Reset()
	}
	// Create a text file to write the output
	f, err := os.Create(filepath.Join(renderConf.OutputFolder, model.TableClassName+".java"))
	if err != nil {
		return err
	}
	defer f.Close()
	// Render java model
	err = t.ExecuteTemplate(f, TEMPLATE_MODEL, model)
	if err != nil {
		return err
	}

	return nil
}

func RenderJava(tables map[string]SqlTable, renderConf conf.RenderConfig) (err error) {
	log.Info("[render] [java] start")
	t, err := initTemplate(renderConf)
	if err != nil {
		return err
	}
	err = initOutputFolder(renderConf)
	if err != nil {
		return err
	}
	// Loop over the tables
	for _, tab := range tables {
		// TODO Exclude pattern filter

		// Init Class Type
		err = initClassFieldType(renderConf, &tab)
		if err != nil {
			log.Errorf("[render] table [%v] init class field type with error: %v", tab.TableName, err)
			continue
		}
		// Render Model File
		err = renderModel(t, renderConf, &tab)
		if err != nil {
			log.Errorf("[render] table [%v] model with error: %v", tab.TableName, err)
			continue
		}
	}
	log.Info("[render] [java] end")
	return nil
}
