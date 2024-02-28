package app

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"chiyou.code/mmc/conf"
	"chiyou.code/mmc/lib/sos"
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
	TEMPLATE_MYBATIS     = "mybatis.tmpl"
	TEMPLATE_MAPPER      = "mapper.tmpl"
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
type ModelView struct {
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

// model.field.tmpl
type ModelFieldView struct {
	FieldComment     string
	FieldAnnotations []string
	FieldType        string
	FieldName        string
}

// mybatis.tmpl
type MybatisView struct {
	// Base Package
	BasePackage string
	// Table Name
	TableName string
	// Table Class Name
	TableClassName string
	// Table Primary Key (Field Name)
	TablePrimaryKeyFieldName string
	// Columns
	Columns []*MybatisColumnView
	// Update Ignore
	UpdateStatementIgnoreColumns []string
	// Delete Ignore
	DeleteStatementIgnoreColumns []string
}

type MybatisColumnView struct {
	// Column Ordinal Position
	ColumnOrdinal int
	// Primary Key
	IsPrimaryKey bool
	// Column Name
	ColumnName string
	// Class Field Name
	ClassFieldName string
}

type MapperView struct {
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
	// Table Primary Key (Field Type)
	TablePrimaryKeyFieldType string
	// Table Primary Key (Field Name)
	TablePrimaryKeyFieldName string
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
		"in": func(value any, collection any) bool {
			switch collection := collection.(type) {
			case []string:
				for _, v := range collection {
					if v == value {
						return true
					}
				}
			case []int:
				for _, v := range collection {
					if v == value {
						return true
					}
				}
			case map[string]interface{}:
				if _, ok := collection[value.(string)]; ok {
					return true
				}
			case map[int]interface{}:
				if _, ok := collection[value.(int)]; ok {
					return true
				}
			}
			return false
		},
		"firstLower": func(value string) string {
			return str.FirstLower(value)
		},
	}).ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Init output folder
func initOutputFolder(renderConf conf.RenderConfig) (err error) {
	if err = sos.CreateFolder(renderConf.OutputFolder); err != nil {
		return err
	}
	return os.RemoveAll(renderConf.OutputFolder)
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
func renderModel(t *template.Template, renderConf conf.RenderConfig, table *SqlTable) (err error) {
	var model ModelView = ModelView{
		Author:         renderConf.Author,
		Version:        renderConf.Version,
		BasePackage:    renderConf.BasePackage,
		TableComment:   table.TableName,
		TableClassName: table.TableClassName,
		CreateTime:     time.Now().Format("2006-01-02 15:04:05"),
		Fields:         make([]string, 0),
	}
	if table.TableComment != nil {
		model.TableComment = *table.TableComment
	}
	// Iterate over the sorted tab
	buf := &bytes.Buffer{}
	for _, column := range table.OrdinalColumns {
		field := ModelFieldView{
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
	modelFileName := filepath.Join(renderConf.OutputFolder, renderConf.ModelSubFolder, model.TableClassName+".java")
	f, err := sos.CreateFile(modelFileName)
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

// Render Mybatis
func renderMybatis(t *template.Template, renderConf conf.RenderConfig, table *SqlTable) (err error) {
	var mybatisView MybatisView = MybatisView{
		BasePackage:                  renderConf.BasePackage,
		TableName:                    table.TableName,
		TableClassName:               table.TableClassName,
		TablePrimaryKeyFieldName:     "",
		Columns:                      make([]*MybatisColumnView, 0),
		UpdateStatementIgnoreColumns: renderConf.UpdateStatementIgnoreColumns,
		DeleteStatementIgnoreColumns: renderConf.DeleteStatementIgnoreColumns,
	}
	for _, column := range table.OrdinalColumns {
		columnView := MybatisColumnView{
			ColumnOrdinal:  column.ColumnOrdinal,
			IsPrimaryKey:   false,
			ColumnName:     column.ColumnName,
			ClassFieldName: column.ClassFieldName,
		}
		if column.ColumnKey != nil && (*column.ColumnKey) == "PRI" {
			columnView.IsPrimaryKey = true
			mybatisView.TablePrimaryKeyFieldName = column.ClassFieldName
		}
		mybatisView.Columns = append(mybatisView.Columns, &columnView)
	}
	// Create a text file to write the output
	modelFileName := filepath.Join(renderConf.OutputFolder, renderConf.DaoSubFolder, mybatisView.TableClassName+"Mapper.xml")
	f, err := sos.CreateFile(modelFileName)
	if err != nil {
		return err
	}
	defer f.Close()
	// Render mybatis
	err = t.ExecuteTemplate(f, TEMPLATE_MYBATIS, mybatisView)
	if err != nil {
		return err
	}
	return nil
}

// Render Mapper
func renderMapper(t *template.Template, renderConf conf.RenderConfig, table *SqlTable) (err error) {
	var mapperView MapperView = MapperView{
		Author:                   renderConf.Author,
		Version:                  renderConf.Version,
		BasePackage:              renderConf.BasePackage,
		TableComment:             table.TableName,
		TableClassName:           table.TableClassName,
		CreateTime:               time.Now().Format("2006-01-02 15:04:05"),
		TablePrimaryKeyFieldType: "",
		TablePrimaryKeyFieldName: "",
	}
	if table.TableComment != nil {
		mapperView.TableComment = *table.TableComment
	}
	// PK
	for _, column := range table.OrdinalColumns {
		if column.ColumnKey != nil && (*column.ColumnKey) == "PRI" {
			mapperView.TablePrimaryKeyFieldType = column.ClassFieldType
			mapperView.TablePrimaryKeyFieldName = column.ClassFieldName
		}
	}
	// Create a text file to write the output
	modelFileName := filepath.Join(renderConf.OutputFolder, "mapper", mapperView.TableClassName+"Mapper.java")
	f, err := sos.CreateFile(modelFileName)
	if err != nil {
		return err
	}
	defer f.Close()
	// Render mybatis
	err = t.ExecuteTemplate(f, TEMPLATE_MAPPER, mapperView)
	if err != nil {
		return err
	}
	return nil
}

// Render dao
// mysql:
//
//	[1] mybatis
//	[2] mapper
func renderDao(t *template.Template, renderConf conf.RenderConfig, table *SqlTable) (err error) {
	switch table.DbType {
	case DB_TYPE_MYSQL:
		// Render mybatis
		err = renderMybatis(t, renderConf, table)
		if err != nil {
			return err
		}
		// Render mapper
		err = renderMapper(t, renderConf, table)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("db type: %v render dao unSupported", table.DbType)
	}

	return nil
}

func RenderJava(tables map[string]*SqlTable, renderConf conf.RenderConfig) (err error) {
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
	for _, table := range tables {
		// TODO Exclude pattern filter
		// TODO 是否继承basemodel处理
		// Init Class Type
		err = initClassFieldType(renderConf, table)
		if err != nil {
			log.Errorf("[render] table [%v] init class field type with error: %v", table.TableName, err)
			continue
		}
		// Render Model File
		err = renderModel(t, renderConf, table)
		if err != nil {
			log.Errorf("[render] table [%v] model with error: %v", table.TableName, err)
			continue
		}
		// Render Dao File
		err = renderDao(t, renderConf, table)
		if err != nil {
			log.Errorf("[render] table [%v] dao with error: %v", table.TableName, err)
			continue
		}
	}
	log.Info("[render] [java] end")
	return nil
}
