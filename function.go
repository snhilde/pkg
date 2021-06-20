// This file contains the information and logic for the Function type.
package pkg

import (
	"bytes"
	"go/doc"
)

// Function holds information about an exported function in a package.
type Function struct {
	// Function's name.
	name string

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
	_ = decl // TODO

	// Extract the parameters.
	in, out := extractParameters(f.Decl.Type, r)

	return Function{
		name:   f.Name,
		input:  in,
		output: out,
	}
}

// Name returns the function's name.
func (f Function) Name() string {
	return f.name
}

// Input returns a list of input parameters sent to this function, or nil on invalid object. If
// there are no input parameters, this returns a slice of size 0..
func (f Function) Input() []Parameter {
	return f.input
}

// Output returns a list of output parameters returned from this function, or nil on invalid object.
// If there are no output parameters, this returns a slice of size 0..
func (f Function) Output() []Parameter {
	return f.output
}
