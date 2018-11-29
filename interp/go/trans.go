package gob

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"text/template"

	"github.com/ClarkGuan/jnif/interp"
)

type goTransform struct {
	packageName string
	outputDir   string
}

func (gt *goTransform) Init(args []string) (string, error) {
	flagSet := flag.NewFlagSet("jnif-go", flag.ContinueOnError)

	flagSet.StringVar(&gt.packageName, "p", "main", "Go package name")
	flagSet.StringVar(&gt.outputDir, "o", ".", "output directory")

	if err := flagSet.Parse(args); err != nil {
		return "", err
	}

	if len(flagSet.Args()) == 0 {
		return "", errors.New("no jar or class file(s) found")
	}

	return flagSet.Arg(0), nil
}

func (gt *goTransform) Transform(methods map[string][]*interp.Method, maxCount int) error {
	obj := map[string]interface{}{
		"classes":     methods,
		"maxCount":    maxCount,
		"packageName": gt.packageName,
	}

	if err := os.MkdirAll(gt.outputDir, 0755); err != nil {
		return err
	}

	headerFile, err := os.OpenFile(filepath.Join(gt.outputDir, "libs.c"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer headerFile.Close()

	goFile, err := os.OpenFile(filepath.Join(gt.outputDir, "libs.go"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer goFile.Close()

	headerFileTpl := template.New("header")
	headerFileTpl.Funcs(template.FuncMap{
		"declaration": declaration,
		"goFuncName":  goFuncName,
	})
	headerFileTpl.Parse(headerTpl)

	goFileTpl := template.New("go")
	goFileTpl.Funcs(template.FuncMap{
		"goFunc": goFunc,
	})
	goFileTpl.Parse(goTpl)

	if err = headerFileTpl.Execute(headerFile, obj); err != nil {
		return err
	}

	if err = goFileTpl.Execute(goFile, obj); err != nil {
		return err
	}

	return nil
}
