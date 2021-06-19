package pkg

import (
	"bytes"
	"go/doc"
)

// Type holds information about an exported type in a package.
type Type struct {
	// Type object from go/doc.
	docType *doc.Type

	// Original declaration in source for this type.
	decl *bytes.Reader

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
		docType:   t,
		decl:      decl,
		functions: functions,
		methods:   methods,
	}
}

// IsValid checks whether or not t is a valid Type object.
func (t Type) IsValid() bool {
	if t.docType == nil {
		return false
	}

	if t.decl == nil {
		return false
	}

	// All checks passed.
	return true
}

// Name returns the type's name.
func (t Type) Name() string {
	if !t.IsValid() {
		return ""
	}

	return t.docType.Name
}

// Functions returns a list of functions that primarily return this type.
func (t Type) Functions() []Function {
	if !t.IsValid() {
		return nil
	}

	return t.functions
}

// Methods returns a list of methods for this type.
func (t Type) Methods() []Method {
	if !t.IsValid() {
		return nil
	}

	return t.methods
}

func (t Type) Exports() []interface{} {
	return nil
}
