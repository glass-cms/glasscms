package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
)

// Simulate FieldMasksForType.
func FieldMasksForType(item any, parent string) map[string]struct{} {
	fields := make(map[string]struct{})

	val := reflect.ValueOf(item)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)

		jsonTag := field.Tag.Get("json")

		// If there's a comma (like `json:"name,omitempty"`), strip out the options
		jsonField := jsonTag
		if jsonField == "" || jsonField == "-" {
			continue
		}

		// Recursively handle nested structs
		if field.Type.Kind() == reflect.Struct {
			nestedFields := FieldMasksForType(reflect.New(field.Type).Elem().Interface(), jsonField)
			for nestedField := range nestedFields {
				fields[nestedField] = struct{}{}
			}
		} else {
			// Add the parent prefix if needed
			if parent != "" {
				jsonField = parent + "." + jsonField
			}
			fields[jsonField] = struct{}{}
		}
	}

	return fields
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("You must provide a type name")
		os.Exit(1)
	}

	typeName := os.Args[1]

	// Get the Go file name from the GOFILE environment variable set by go generate
	goFile := os.Getenv("GOFILE")
	if goFile == "" {
		fmt.Println("GOFILE environment variable is not set")
		os.Exit(1)
	}

	// Parse the Go file that contains the type declaration
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, goFile, nil, parser.AllErrors)
	if err != nil {
		fmt.Println("Error parsing Go file:", err)
		os.Exit(1)
	}

	// Find the struct type by name
	var targetStruct *ast.StructType
	ast.Inspect(node, func(n ast.Node) bool {
		if ts, ok := n.(*ast.TypeSpec); ok {
			if ts.Name.Name == typeName {
				if structType, isStruct := ts.Type.(*ast.StructType); isStruct {
					targetStruct = structType
					return false // stop after we find the correct type
				}
			}
		}
		return true
	})

	if targetStruct == nil {
		fmt.Printf("Type %s not found in file\n", typeName)
		os.Exit(1)
	}

	// Dynamically create a value of the given type (using reflection)
	// Here we assume a dynamic way to generate a reflection value for the type (manual mapping can be done)
	// This part may need to be extended to fully work for your case.
	dynamicType := reflect.New(reflect.StructOf([]reflect.StructField{})).Elem().Interface()
	fmt.Println(dynamicType)
	// Get field masks
	fields := FieldMasksForType(dynamicType, "")
	fmt.Println("Do I even get here?")
	// Print out the field masks
	for field := range fields {
		fmt.Println(field)
	}
}
