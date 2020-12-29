# JSONSchema validator for JSON and YAML files

## Features

* Validate JSON and YAML files with your JSON Schema file

## Usage

*Note:* Right now both the schema and the file to be validated have to be located on the file system.

### Binary

Set the schema as first argument and the file you want to validate as second argument.

Alternatively you may also use the environment variables `SCHEMA_PATH` and `FILE_PATH`.

```
bin/jsonschema-validator ./schema/project-definition.schema.json ./examples/valid/enum1.yml

SCHEMA_PATH=./schema/project-definition.schema.json FILE_PATH=./examples/valid/enum1.json bin/jsonschema-validator
```

### GoLang API

```go
package main

import validator "github.com/saitho/jsonschema-validator/validator"

// There is also ValidateFileWithInput if you want to validate with the schema data itself (byte[])
result, err := validator.ValidateFile(filePath, schemaPath) // result is gojsonschema.Result
errorMessage := validator.ShouldValidate(result)
if len(errorMessage) == 0 {
	_, err = fmt.Fprintf(os.Stdout, "The file is valid.\n")
} else {
	_, err = fmt.Fprintf(os.Stderr, errorMessage + "\n")
	if err != nil {
		panic(err.Error())
	}
	os.Exit(1)
}
```

### Docker container

You can find the official Docker image on [Docker Hub](https://hub.docker.com/r/saitho/jsonschema-validator).
The binary on the Docker container is located at `/bin/validator`.

Example call (currently broken! Drone works fine)
```
docker run --rm -u=$(id -u):$(id -g) -v=$(pwd):/app saitho/jsonschema-validator ./schema/project-definition.schema.json ./examples/valid/enum1.yml
```

### DroneCI

```yaml
---
kind: pipeline
name: Linter
type: docker

steps:
  - name: Validate providers.yml
    image: saitho/jsonschema-validator
    commands:
      - /bin/validator ./tests/schema.json ./config/configuration.yml
 ```

## Development

### Run Tests

```shell script
go test
```

### Build binary

```shell script
go build -o bin/jsonschema-validator
```
