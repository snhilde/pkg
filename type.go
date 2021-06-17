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

func (t Type) Functions() []Function {
	return nil
}

func (t Type) Methods() []Method {
	return nil
}

func (t Type) Exports() []interface{} {
	return nil
}

// Method holds information about a type's method.
type Method struct {
}
