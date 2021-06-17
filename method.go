package pkg

import (
	"go/doc"
)

// Method holds information about a type's method.
type Method struct {
	// Func object from go/doc.
	docMethod *doc.Func
}

// isValid checks whether or not m is a valid Method object.
func (m Method) isValid() bool {
	if m == (Method{}) {
		return false
	}

	if m.docMethod == nil {
		return false
	}

	// All checks passed.
	return true
}

// Name returns the method's name.
func (m Method) Name() string {
	if !m.isValid() {
		return ""
	}

	return m.docMethod.Name
}
