package ir

import (
	"os"
	"path/filepath"
	"text/template"
)

func Print(path string, packageName string, dir string) error {
	methods, max, err := Parse(path)
	if err != nil {
		return err
	}

	if len(packageName) == 0 {
		packageName = "main"
	}

	if len(dir) == 0 {
		dir = "."
	}

	obj := map[string]interface{}{
		"classes":     methods,
		"maxCount":    max,
		"packageName": packageName,
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	headerFile, err := os.OpenFile(filepath.Join(dir, "libs.c"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer headerFile.Close()

	goFile, err := os.OpenFile(filepath.Join(dir, "libs.go"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer goFile.Close()

	headerFileTpl := template.Must(template.New("header").Parse(headerTpl))
	goFileTpl := template.Must(template.New("go").Parse(goTpl))

	if err = headerFileTpl.Execute(headerFile, obj); err != nil {
		return err
	}

	if err = goFileTpl.Execute(goFile, obj); err != nil {
		return err
	}

	return nil
}
