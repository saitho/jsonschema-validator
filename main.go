package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"
)

// isInternalError determines if a given error message is related to the schema requirements itself
func isInternalError(errorType string) bool {
	switch errorType {
	case
		"condition_else",
		"condition_then",
		"number_any_of",
		"number_one_of",
		"number_all_of",
		"number_not":
		return true
	default:
		return false
	}
}

// ValidateFile validates the contents of filePath with the schema
func ValidateFile(filePath string, schemaPath string) (*gojsonschema.Result, error) {
	var err error

	matchRegex, err := regexp.MatchString("\\w+://", schemaPath)
	if !matchRegex {
		schemaPath = "file://" + schemaPath
	}
	schemaLoader := gojsonschema.NewReferenceLoader(schemaPath)

	configJson, err := loadJsonFile(filePath)
	if err != nil {
		return nil, err
	}
	documentLoader := gojsonschema.NewBytesLoader(configJson)

	return gojsonschema.Validate(schemaLoader, documentLoader)
}

// Loads the contents of a given JSON file.
// If the file is a YAML/YML file, it will be converted to JSON
func loadJsonFile(filePath string) ([]byte, error) {
	var err error

	configData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(filePath, ".yml") || strings.HasSuffix(filePath, ".yaml") {
		return yaml.YAMLToJSON(configData)
	} else if !strings.HasSuffix(filePath, ".json") {
		return nil, fmt.Errorf("unknown suffix. allowed values: .json, .yml .yaml")
	}
	return configData, nil
}

// ShouldValidate validates result
// signature uses interface{} and unused parmeter because this method is also used in tests with Convey
func ShouldValidate(actual interface{}, _ ...interface{}) string {
	result := actual.(*gojsonschema.Result)
	if result.Valid() == true {
		return ""
	}
	errorMessage := fmt.Sprintf("INVALID. See errors:\n")

	for _, desc := range result.Errors() {
		if isInternalError(desc.Type()) {
			continue
		}
		errorMessage += fmt.Sprintf("- %s\n", desc)
	}
	return errorMessage
}

func main() {
	if len(os.Args) < 3 {
		panic("Missing argument. Set the schema file path as first command argument and the file you want to validate as second argument.")
	}
	var schemaPath = os.Args[1]
	var filePath = os.Args[2]

	result, err := ValidateFile(filePath, schemaPath)

	if err != nil {
		panic(err.Error())
	}

	errorMessage := ShouldValidate(result)
	if len(errorMessage) == 0 {
		_, err = fmt.Fprintln(os.Stdout, "VALID")
	} else {
		_, err = fmt.Fprintln(os.Stderr, errorMessage)
		os.Exit(1)
	}

	if err != nil {
		panic(err.Error())
	}
	os.Exit(0)
}
