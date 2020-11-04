package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func collectYamlFiles(folder string) []string {
	var files []string
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".json") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

func getSchemaFile() string {
	_, filename, _, _ := runtime.Caller(0)
	binDir, _ := filepath.Abs(filepath.Dir(filename))
	return filepath.Join(binDir, "..", "schema", "project-definition.schema.json")
}

func TestValidDefinitions(t *testing.T) {
	Convey("test valid definitions", t, func() {
		for _, file := range collectYamlFiles("../examples/valid") {
			Convey(fmt.Sprintf("file %s should validate", file), func() {
				result, err := ValidateFile(file, getSchemaFile())
				So(err, ShouldBeNil)
				So(result, ShouldValidate)
			})
		}
	})
}

func TestFolderTraversion(t *testing.T) {
	Convey("test folder traversion", t, func() {
		Convey("all files validate", func() {
			results, err := ValidateDirectory("examples/valid", getSchemaFile())
			So(err, ShouldBeNil)
			for _, result := range results {
				So(result, ShouldValidate)
			}
		})

		Convey("validate all files and return error", func() {
			results, err := ValidateDirectory("examples/invalid", getSchemaFile())
			So(err, ShouldBeNil)
			for _, result := range results {
				So(result, ShouldNotValidate)
			}
		})
	})
}

func TestInvalidDefinitions(t *testing.T) {
	Convey("test invalid definitions", t, func() {
		for _, file := range collectYamlFiles("examples/invalid") {
			Convey(fmt.Sprintf("file %s should not validate", file), func() {
				result, err := ValidateFile(file, getSchemaFile())
				So(err, ShouldBeNil)
				So(result, ShouldNotValidate)
			})
		}
	})
}
