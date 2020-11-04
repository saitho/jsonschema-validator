package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

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
func ValidateFile(filePath string) (*gojsonschema.Result, error) {
	var err error

	binDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + filepath.Join(binDir, "..", "schema", "project-definition.schema.json"))

	configData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// YAML to JSON
	configJson, err := yaml.YAMLToJSON(configData)
	if err != nil {
		return nil, err
	}
	documentLoader := gojsonschema.NewBytesLoader(configJson)

	return gojsonschema.Validate(schemaLoader, documentLoader)
}

// ShouldValidate validates result
// signature uses interface{} and unused paremter because this method is also used in tests with Convey
func ShouldValidate(actual interface{}, _ ...interface{}) string {
	result := actual.(*gojsonschema.Result)
	if result.Valid() == true {
		return ""
	}
	errorMessage := fmt.Sprintf("The project definition is not valid. see errors:\n")

	for _, desc := range result.Errors() {
		if isInternalError(desc.Type()) {
			continue
		}
		errorMessage += fmt.Sprintf("- %s\n", desc)
	}
	return errorMessage
}

func main() {
	if len(os.Args) < 2 {
		panic("Missing argument. Set the file out want to validate as first command argument.")
	}
	var filePath = os.Args[1]

	result, err := ValidateFile(filePath)

	if err != nil {
		panic(err.Error())
	}

	errorMessage := ShouldValidate(result)
	if len(errorMessage) == 0 {
		_, err = fmt.Fprintln(os.Stderr, "The project definition is valid")
	} else {
		_, err = fmt.Fprintln(os.Stderr, errorMessage)
		os.Exit(1)
	}

	if err != nil {
		panic(err.Error())
	}
	os.Exit(0)
}
