package pkg

import (
	"bytes"
	"go/doc"
)

// Type holds information about an exported type in a package.
type Type struct {
	// Name of this type.
	name string

	// Functions in the package that primarily return this type.
	functions []Function

	// Methods for this type.
	methods []Method
}

// newType builds a new Type object based on go/doc's Type.
func newType(t *doc.Type, r *bytes.Reader) Type {
	if t == nil || r == nil {
		return Type{}
	}
	// Read out the source declaration.
	start, end := t.Decl.Pos()-1, t.Decl.End()-1 // -1 to index properly
	decl := extractSource(r, start, end)
	_ = decl // TODO

	// Make a list of functions for this type.
	functions := make([]Function, len(t.Funcs))
	for i, f := range t.Funcs {
		functions[i] = newFunction(f, r)
	}

	// Make a list of methods for this type.
	methods := make([]Method, len(t.Methods))
	for i, m := range t.Methods {
		methods[i] = newMethod(m, r)
	}

	return Type{
		name:      t.Name,
		functions: functions,
		methods:   methods,
	}
}

// Name returns the type's name.
func (t Type) Name() string {
	return t.name
}

// Functions returns a list of functions that primarily return this type.
func (t Type) Functions() []Function {
	return t.functions
}

// Methods returns a list of methods for this type.
func (t Type) Methods() []Method {
	return t.methods
}

func (t Type) Exports() []interface{} {
	return nil
}
