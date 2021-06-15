package pkg

import (
	"go/doc"
)

// Type holds information about an exported type in a package.
type Type struct {
	// Type object from go/doc.
	docType *doc.Type
}

func (t Type) Functions() []Function {
}

func (t Type) Methods() []Method {
}

func (t Type) Exports() []interface{} {
}

// Method holds information about a type's method.
type Method struct {
}
