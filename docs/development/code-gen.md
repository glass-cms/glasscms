# Code Generation

This project uses [TypeSpec](https://typespec.io/) for describing the API of the application. The TypeSpec is compiled in an OpenAPI 3.0 specification file, which is then used to generate the client and server code.

## Generating the OpenAPI specification

The TypeSpec files are located in the `typespec` directory. To compile the OpenAPI specification, run the following command:

```bash
make typespec
```

The OpenAPI specification will be generated in the `openapi` directory.

## Generating the client and server code

The client and server code is generated using the [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) module. To generate the code, run the following command:

```bash
go generate ./...
```
