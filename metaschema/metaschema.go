package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"moonspace/model"
	"os"
	"strings"

	"github.com/invopop/jsonschema"
)

const location = "../schemas/"
const prefix = "#/$defs/"

func GenerateMetaschema() {
	err := os.RemoveAll(location)
	if err != nil {
		panic(err)
	}

	err = os.Mkdir(location, os.ModePerm)
	if err != nil {
		panic(fmt.Errorf("couldn't create a directory to place schemas in; %v", err))
	}

	schemas := map[string]*jsonschema.Schema{
		"category": jsonschema.Reflect(&model.Category{}),
		"product":  jsonschema.Reflect(&model.Product{}),
		"cart":     jsonschema.Reflect(&model.Cart{}),
		"order":    jsonschema.Reflect(&model.Order{}),
	}

	for k, s := range schemas {
		for _, d := range s.Definitions {
			ref := strings.ToLower(strings.ReplaceAll(s.Ref, prefix, ""))
			if ref == k {
				saveSchema(k, d)
				break
			}
		}

	}
}

func saveSchema(schemaName string, d *jsonschema.Schema) {
	bytes, err := d.Properties.MarshalJSON()
	if err != nil {
		panic(err)
	}

	if len(bytes) > 0 {
		filename := location + schemaName + ".json"
		err = ioutil.WriteFile(filename, bytes, fs.ModeAppend)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	GenerateMetaschema()
}
