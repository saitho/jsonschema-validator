package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// copies schema file from code folder to the folder where the test is executed (usually /tmp/go-build[number]/)
func setupSchemaFile() {
	_, filename, _, _ := runtime.Caller(0)
	binDir, _ := filepath.Abs(filepath.Dir(filename))
	var schemaFile = filepath.Join(binDir, "..", "schema", "project-definition.schema.json")

	pwdBinDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	var pwdSchemaFolder = filepath.Join(pwdBinDir, "..", "schema")
	var pwdSchemaFile = filepath.Join(pwdSchemaFolder, "project-definition.schema.json")

	_ = os.MkdirAll(pwdSchemaFolder, os.ModePerm)

	input, err := ioutil.ReadFile(schemaFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile(pwdSchemaFile, input, 0644)
	if err != nil {
		fmt.Println("Error creating", pwdSchemaFile)
		fmt.Println(err)
		return
	}
}

func collectYamlFiles(folder string) []string {
	var files []string
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

func TestValidDdefinitions(t *testing.T) {
	setupSchemaFile()
	Convey("test valid definitions", t, func() {
		for _, file := range collectYamlFiles("../examples/valid") {
			Convey(fmt.Sprintf("file %s should validate", file), func() {
				result, err := ValidateFile(file)
				So(err, ShouldBeNil)
				So(result, ShouldValidate)
			})
		}
	})
}

func TestInvalidDefinitions(t *testing.T) {
	setupSchemaFile()
	Convey("test invalid definitions", t, func() {
		for _, file := range collectYamlFiles("../examples/invalid") {
			Convey(fmt.Sprintf("file %s should not validate", file), func() {
				result, err := ValidateFile(file)
				So(err, ShouldBeNil)
				So(result, ShouldNotValidate)
			})
		}
	})
}

func ShouldNotValidate(actual interface{}, _ ...interface{}) string {
	result := ShouldValidate(actual)
	if result == "" { // file validated
		return "File validated (it should not)"
	} else {
		return ""
	}
}
