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

// Render model object
func renderModel(t *template.Template, javaConf conf.RenderConfig, tab SqlTable) (err error) {
	rep := strings.NewReplacer(javaConf.IgnorePrefix, "", javaConf.IgnoreSuffix, "")
	var model Model = Model{
		BasePackage:    javaConf.BasePackage,
		TableComment:   tab.TableName,
		TableClassName: str.UnderscoreToCamel(rep.Replace(tab.TableName)),
		CreateTime:     time.Now().Format("2006-01-02 15:04:05"),
		Fields:         make([]string, len(tab.Columns)),
	}
	if tab.TableComment != nil {
		model.TableComment = *tab.TableComment
	}

	buf := &bytes.Buffer{}
	for colName, column := range tab.Columns {
		field := ModelField{
			FieldComment:    *column.ColumnComment,
			FieldAnnotation: "@Min(value = 0)",
			FieldType:       "Integer",
			FieldName:       colName,
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
		err = renderModel(t, javaConf, tab)
		if err != nil {
			return err
		}
	}

	return nil
}
