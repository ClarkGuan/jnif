package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ClarkGuan/jnif/ir"
)

func main() {
	var packageName string
	var outputDir string

	flag.StringVar(&packageName, "p", "main", "Go package name")
	flag.StringVar(&outputDir, "o", ".", "output directory")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Fprintln(os.Stderr, "no jar or class file(s) found")
		flag.Usage()
		os.Exit(1)
	}

	if err := ir.Print(flag.Arg(0), packageName, outputDir); err != nil {
		log.Fatalln(err)
	}
}
