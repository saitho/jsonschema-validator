# JSONSchema validator for JSON and YAML files

## Features

* Validate JSON and YAML files with your JSON Schema file

## Usage

*Note:* Right now both the schema and the file to be validated have to be located on the file system.

### Binary

```
bin/jsonschema-validator ./schema/project-definition.schema.json ./examples/valid/enum1.yml

bin/jsonschema-validator ./schema/project-definition.schema.json ./examples/valid/enum1.json
```

### GoLang API

### Docker container

You can find the official Docker image on [Docker Hub](https://hub.docker.com/r/saitho/jsonschema-validator).
The binary on the Docker container is located at `/bin/validator`.

```
docker run --rm -u=$(id -u):$(id -g) -v=$(pwd):/app saitho/jsonschema-validator ./schema/project-definition.schema.json ./examples/valid/enum1.yml
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
