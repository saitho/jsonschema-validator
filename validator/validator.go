package validator

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ValidateDirectory(directoryPath string, schemaPath string) ([]*gojsonschema.Result, error) {
	var results []*gojsonschema.Result

	var validationResult *gojsonschema.Result
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".json") {
			validationResult, err = ValidateFile(path, schemaPath)
			if err != nil {
				return err
			}
			results = append(results, validationResult)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return results, nil
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
