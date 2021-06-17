package pkg

import (
	"go/doc"
)

// Type holds information about an exported type in a package.
type Type struct {
	// Type object from go/doc.
	docType *doc.Type
}

// isValid checks whether or not t is a valid Type object.
func (t Type) isValid() bool {
	if t == (Type{}) {
		return false
	}

	if t.docType == nil {
		return false
	}

	// All checks passed.
	return true
}

// Name returns the type's name.
func (t Type) Name() string {
	if !t.isValid() {
		return ""
	}

	return t.docType.Name
}

// Functions returns a list of functions that primarily return this type.
func (t Type) Functions() []Function {
	if !t.isValid() {
		return nil
	}

	// Wrap every go/doc Func in our own Function.
	functions := make([]Function, len(t.docType.Funcs))
	for i, f := range t.docType.Funcs {
		functions[i] = Function{
			docFunc: f,
		}
	}

	return functions
}

// Methods returns a list of methods for this type.
func (t Type) Methods() []Method {
	if !t.isValid() {
		return nil
	}

	// Wrap every go/doc Func in our own Method.
	methods := make([]Method, len(t.docType.Methods))
	for i, m := range t.docType.Methods {
		methods[i] = Method{
			docMethod: m,
		}
	}

	return methods
}

func (t Type) Exports() []interface{} {
	return nil
}
