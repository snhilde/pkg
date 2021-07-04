// This file contains the information and logic for the Type type.
package pkg

import (
	"go/ast"
	"go/doc"
	"go/token"
)

// Type holds information about an exported type in a package.
type Type struct {
	// Name of this type.
	name string

	// Comments for this type.
	comments string

	// Name of this type's underlying type.
	typeName string

	// Original declaration in source for this type.
	source string

	// Functions in the package that primarily return this type.
	functions []Function

	// Methods for this type.
	methods []Method
}

// newType builds a new Type object based on go/doc's Type.
func newType(t *doc.Type, fset *token.FileSet) Type {
	if t == nil {
		return Type{}
	}

	// Read out the source declaration (exported fields only).
	source := extractSource(t.Decl, fset)

	// Extract the underlying type.
	typeName := extractType(t, fset)

	// Make a list of functions for this type.
	functions := make([]Function, len(t.Funcs))
	for i, f := range t.Funcs {
		functions[i] = newFunction(f, fset)
	}

	// Make a list of methods for this type.
	methods := make([]Method, len(t.Methods))
	for i, m := range t.Methods {
		methods[i] = newMethod(m, fset)
	}

	return Type{
		name:      t.Name,
		comments:  t.Doc,
		typeName:  typeName,
		source:    source,
		functions: functions,
		methods:   methods,
	}
}

// extractType extracts the underlying type name for this type.
func extractType(t *doc.Type, fset *token.FileSet) string {
	if t == nil || t.Decl == nil || len(t.Decl.Specs) == 0 {
		return ""
	}

	ts := t.Decl.Specs[0].(*ast.TypeSpec)
	switch ts.Type.(type) {
	case *ast.StructType:
		return "struct"
	case *ast.InterfaceType:
		return "interface"
	default:
		return extractSource(ts.Type, fset)
	}
}

// Name returns the type's name.
func (t Type) Name() string {
	return t.name
}

// Comments returns the documentation for this type with pkg's formatting applied.
func (t Type) Comments(width int) string {
	return formatComments(t.comments, width)
}

// Type returns the name of the type's underlying type, like "struct" or "map[string]chan int".
func (t Type) Type() string {
	return t.typeName
}

// Source returns the source declaration for this type.
func (t Type) Source() string {
	return t.source
}

// Functions returns a list of functions that primarily return this type.
func (t Type) Functions() []Function {
	return append([]Function{}, t.functions...)
}

// Methods returns a list of methods for this type.
func (t Type) Methods() []Method {
	return append([]Method{}, t.methods...)
}
