package app

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"chiyou.code/mmc/conf"
)

// Person is a struct that holds some information about a person
type Person struct {
	Name    string
	Age     int
	Friends []string
}

func RenderJava(tables map[string]SqlTable, javaConf conf.JavaConfig) (err error) {
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
		return err
	}
	t, err := template.ParseFiles(files...)
	if err != nil {
		return err
	}

	// Check if the directory exists
	_, err = os.Stat(javaConf.OutputFolder)
	if errors.Is(err, os.ErrNotExist) {
		// Create the directory if it does not exist
		err = os.MkdirAll(javaConf.OutputFolder, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	// Create a text file to write the output
	f, err := os.Create(filepath.Join(javaConf.OutputFolder, "output.txt"))
	if err != nil {
		return err
	}
	defer f.Close()

	// Create a slice of Person objects
	people := []Person{
		{Name: "Alice", Age: 25, Friends: []string{"Bob", "Carol"}},
		{Name: "Bob", Age: 26, Friends: []string{"Alice", "Dave"}},
		{Name: "Carol", Age: 27, Friends: []string{"Alice", "Eve"}},
	}
	// Execute the template for each person and write to the file
	for _, p := range people {
		err = t.ExecuteTemplate(f, "model.tmpl", p)
		if err != nil {
			return err
		}
	}

	// Print a success message
	log.Println("Output written to output.txt")

	return nil
}
