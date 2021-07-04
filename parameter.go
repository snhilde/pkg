// This file contains the logic for extracting and using parameters to/from functions.
package pkg

import (
	"go/ast"
	"go/token"
	"strings"
)

// Parameter represents a parameter to a function, either as input (argument) or output (result).
type Parameter struct {
	// Name of the parameter.
	name string

	// Name of the parameter's type.
	typeName string
}

// newParameters extracts all parameters for the given list of fields.
func newParameters(list *ast.FieldList, fset *token.FileSet) []Parameter {
	if list == nil {
		return nil
	}

	params := make([]Parameter, 0)
	for _, p := range list.List {
		// Read out the type for this parameter(s).
		typeName := extractSource(p.Type, fset)
		typeName = strings.TrimSpace(typeName)

		if p.Names == nil {
			// Unnamed parameter (probably a return).
			params = append(params, Parameter{
				name:     "",
				typeName: typeName,
			})
		} else {
			// Get the names of all grouped parameters of this type and add the members to the list.
			for _, name := range p.Names {
				params = append(params, Parameter{
					name:     strings.TrimSpace(name.Name),
					typeName: typeName,
				})
			}
		}
	}

	return params
}

// Name returns the parameter's name. Might be empty for output parameters.
func (p Parameter) Name() string {
	return p.name
}

// Type returns the parameter's type, like "[]byte" or "*os.File".
func (p Parameter) Type() string {
	return p.typeName
}

// Pointer reports whether or not this parameter is a pointer. If the parameter is a slice of
// pointers, this returns false.
func (p Parameter) Pointer() bool {
	// Remove the ... prefix for variadic parameters.
	s := strings.TrimPrefix(p.typeName, "...")

	return strings.HasPrefix(s, "*")
}

// String returns the string representation of this type.
func (p Parameter) String() string {
	s := p.name + " " + p.typeName

	return strings.TrimSpace(s)
}
