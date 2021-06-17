package pkg

import (
	"go/doc"
)

// Function holds information about an exported function in a package.
type Function struct {
	docFunc *doc.Func
}

// isValid checks whether or not f is a valid Function object.
func (f Function) isValid() bool {
	if f == (Function{}) {
		return false
	}

	if f.docFunc == nil {
		return false
	}

	// All checks passed.
	return true
}

// Name returns the function's name.
func (f Function) Name() string {
	if !f.isValid() {
		return ""
	}

	return f.docFunc.Name
}
