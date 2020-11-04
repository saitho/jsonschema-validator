package main

import (
	"fmt"
	"os"

	"github.com/saitho/jsonschema-validator/validator"
	"github.com/xeipuuv/gojsonschema"
)

func evaluateResult(result *gojsonschema.Result) bool {
	var err error

	errorMessage := validator.ShouldValidate(result)
	if len(errorMessage) == 0 {
		_, err = fmt.Fprintln(os.Stdout, "VALID")
	} else {
		_, err = fmt.Fprintln(os.Stderr, errorMessage)
		return false
	}

	if err != nil {
		panic(err.Error())
	}
	return true
}

func main() {
	if len(os.Args) < 3 {
		panic("Missing argument. Set the schema file path as first command argument and the file you want to validate as second argument.")
	}
	var schemaPath = os.Args[1]
	var filePath = os.Args[2]

	hasErrors := false
	stat, _ := os.Stat(filePath)
	if stat.IsDir() {
		results, err := validator.ValidateDirectory(filePath, schemaPath)
		if err != nil {
			panic(err.Error())
		}
		for _, result := range results {
			if !evaluateResult(result) {
				hasErrors = true
			}
		}
	} else {
		result, err := validator.ValidateFile(filePath, schemaPath)
		if err != nil {
			panic(err.Error())
		}
		hasErrors = !evaluateResult(result)
	}

	if hasErrors {
		os.Exit(1)
	}
	os.Exit(0)
}
