// This file contains the information and logic for the Function type.
package pkg

import (
	"bytes"
	"go/doc"
)

// Function holds information about an exported function in a package.
type Function struct {
	// Func object from go/doc.
	docFunc *doc.Func

	// Function's name.
	name string

	// Original declaration in source for this function.
	decl *bytes.Reader

	// Input parameters.
	input []Parameter

	// Output parameters.
	output []Parameter
}

// newFunction builds a new Function object based on go/doc's Func.
func newFunction(f *doc.Func, r *bytes.Reader) Function {
	if f == nil || r == nil {
		return Function{}
	}
	// Read out the source declaration.
	start, end := f.Decl.Type.Pos()-1, f.Decl.Type.End()-1 // -1 to index properly
	decl := extractSource(r, start, end)

	// Extract the parameters.
	in, out := extractParameters(f.Decl.Type, r)

	return Function{
		docFunc: f,
		name:    f.Name,
		decl:    decl,
		input:   in,
		output:  out,
	}
}

// IsValid checks whether or not f is a valid Function object.
func (f Function) IsValid() bool {
	if f.docFunc == nil {
		return false
	}

	if f.name == "" {
		return false
	}

	if f.decl.Len() == 0 {
		return false
	}

	// All checks passed.
	return true
}

// Name returns the function's name.
func (f Function) Name() string {
	if !f.IsValid() {
		return ""
	}

	return f.name
}

// Input returns a list of input parameters sent to this function, or nil on invalid object. If
// there are no input parameters, this returns a slice of size 0..
func (f Function) Input() []Parameter {
	if !f.IsValid() {
		return nil
	}

	return f.input
}

// Output returns a list of output parameters returned from this function, or nil on invalid object.
// If there are no output parameters, this returns a slice of size 0..
func (f Function) Output() []Parameter {
	if !f.IsValid() {
		return nil
	}

	return f.output
}
