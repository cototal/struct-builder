package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type PackageStructs struct {
	Package string       `json:"package"`
	Imports []string     `json:"imports"`
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
	for _, packageStructList := range infoMap {
		var structBuilder strings.Builder
		structBuilder.WriteString(fmt.Sprintf("package %s\n\n", packageStructList.Package))
		if len(packageStructList.Imports) > 0 {
			structBuilder.WriteString("imports (\n")
			for _, importPackage := range packageStructList.Imports {
				structBuilder.WriteString(fmt.Sprintf("\t\"%s\"\n", importPackage))
			}
			structBuilder.WriteString(")")
		}

		for _, structInfo := range packageStructList.Structs {
			structBuilder.WriteString(fmt.Sprintf("\n\ntype %s struct {\n", structInfo.Name))
			for field, dataType := range structInfo.Fields {
				structBuilder.WriteString(fmt.Sprintf("\t%s %s ", field, dataType))
				if len(structInfo.Tags) == 0 {
					continue
				}
				var tagBuilder strings.Builder
				for _, tagType := range structInfo.Tags {
					tagValue := fmt.Sprintf("%s:\"%s\" ", tagType, strings.ToLower(field))
					switch tagType {
					case "json":
						override, ok := structInfo.JSON[field]
						if ok {
							tagValue = fmt.Sprintf("json:\"%s\" ", override)
						}
					}
					tagBuilder.WriteString(tagValue)
				}
				tagString := strings.TrimSpace(tagBuilder.String())
				structBuilder.WriteString(fmt.Sprintf("`%s`\n", tagString))
			}
			structBuilder.WriteString("}\n")
		}

		packageDir := workDir
		if packageStructList.Package != "main" {
			packageDir = filepath.Join(workDir, packageStructList.Package)
		}
		if err := os.MkdirAll(packageDir, os.ModePerm); err != nil {
			log.Fatal(err)
		}
		structFile := filepath.Join(packageDir, "structs.tmp")
		err = os.WriteFile(structFile, []byte(structBuilder.String()), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
}
