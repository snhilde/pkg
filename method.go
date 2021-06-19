// This file contains the information and logic for the Method type.
package pkg

import (
	"bytes"
	"go/doc"
	"strings"
)

// Method holds information about a type's method.
type Method struct {
	// Func object from go/doc.
	docFunc *doc.Func

	// Method's name.
	name string

	// Original declaration in source for this method.
	decl *bytes.Reader

	// Input parameters.
	input []Parameter

	// Output parameters.
	output []Parameter

	// Whether or not the receiver of this method is a pointer.
	pointer bool
}

// newMethod builds a new Method object based on go/doc's Func.
func newMethod(f *doc.Func, r *bytes.Reader) Method {
	if f == nil || r == nil {
		return Method{}
	}
	// Read out the source declaration.
	start, end := f.Decl.Type.Pos()-1, f.Decl.Type.End()-1 // -1 to index properly
	decl := extractSource(r, start, end)

	// Extract the parameters.
	in, out := extractParameters(f.Decl.Type, r)

	return Method{
		docFunc: f,
		name:    f.Name,
		decl:    decl,
		input:   in,
		output:  out,
		pointer: strings.HasPrefix(f.Recv, "*"),
	}
}

// IsValid checks whether or not m is a valid Method object.
func (m Method) IsValid() bool {
	if m.docFunc == nil {
		return false
	}

	if m.name == "" {
		return false
	}

	if m.decl.Len() == 0 {
		return false
	}

	// All checks passed.
	return true
}

// Name returns the method's name.
func (m Method) Name() string {
	if !m.IsValid() {
		return ""
	}

	return m.name
}

// Input returns a list of input parameters sent to this method, or nil on invalid object. If
// there are no input parameters, this returns a slice of size 0.. The list does not include the
// method receiver.
func (m Method) Input() []Parameter {
	if !m.IsValid() {
		return nil
	}

	return m.input
}

// Output returns a list of output parameters returned from this method, or nil on invalid object.
// If there are no output parameters, this returns a slice of size 0..
func (m Method) Output() []Parameter {
	if !m.IsValid() {
		return nil
	}

	return m.output
}
