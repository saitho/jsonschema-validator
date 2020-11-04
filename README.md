# JSONSchema validator for JSON and YAML files

## Features

* Validate JSON and YAML files with your JSON Schema file

## Usage

*Note:* Right now both the schema and the file to be validated have to be located on the file system.

### Binary

```
bin/jsonschema-validator ./schema/project-definition.schema.json ./examples/valid/enum1.yml
```

### GoLang API

### Docker container

## Development

### Run Tests

```shell script
go test
```

### Build binary

```shell script
go build -o bin/jsonschema-validator
```
