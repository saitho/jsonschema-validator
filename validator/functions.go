package validator

import (
	"fmt"
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

func ShouldNotValidate(actual interface{}, _ ...interface{}) string {
	result := ShouldValidate(actual)
	if result == "" { // file validated
		return "File validated (it should not)"
	} else {
		return ""
	}
}
