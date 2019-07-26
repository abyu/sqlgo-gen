package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/abyu/sqlxx/db_template"
	"github.com/abyu/sqlxx/generator"
)

var (
	srcPath  = flag.String("source", "", "the source directory that has all the entities to generate stores for")
	destPath = flag.String("destination", "", "the destination directory to save the generated files")
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "\tsqlxx [flags]:\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = Usage
	flag.Parse()
	if len(*srcPath) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	if len(*destPath) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	typeInfo := generator.LoadPackage(*srcPath)

	if _, err := os.Stat(*destPath); os.IsNotExist(err) {
		os.Mkdir(*destPath, os.ModePerm)
	}

	repoTemplate := db_template.NewMysqlTemplate()

	for _, info := range typeInfo {
		generator.GenerateStruct(*destPath, info.Name, info, repoTemplate)
	}
}
