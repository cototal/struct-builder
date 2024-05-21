package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type PackageStructs struct {
	Package string       `json:"package"`
	Structs []StructInfo `json:"structs"`
}

type StructInfo struct {
	Name   string            `json:"name"`
	Tags   []string          `json:"tags"`
	Fields map[string]string `json:"fields"`
	JSON   map[string]string `json:"json"`
	XML    map[string]string `json:"xml"`
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Usage...")
		os.Exit(0)
	}

	fileName := args[0]
	file, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	var infoMap []PackageStructs
	err = json.Unmarshal(file, &infoMap)
	if err != nil {
		log.Fatal(err)
	}

	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	for _, structItem := range infoMap {
		packageDir := workDir
		if structItem.Package != "main" {
			packageDir = filepath.Join(workDir, structItem.Package)
		}
		if err := os.MkdirAll(packageDir, os.ModePerm); err != nil {
			log.Fatal(err)
		}
		structFile := filepath.Join(packageDir, "structs.tmp")
		err = os.WriteFile(structFile, []byte(structItem.Package), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
}
