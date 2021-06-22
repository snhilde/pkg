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

	// Comments for this function.
	comments string

	// Input parameters.
	inputs []Parameter

	// Output parameters.
	outputs []Parameter
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
		name:     f.Name,
		comments: f.Doc,
		inputs:   in,
		outputs:  out,
	}
}

// Name returns the function's name.
func (f Function) Name() string {
	return f.name
}

// Comments returns the documentation for this function with pkg's formatting applied.
func (f Function) Comments(width int) string {
	return formatComments(f.comments, width)
}

// Inputs returns a list of input parameters sent to this function, or nil on invalid object. If
// there are no input parameters, this returns a slice of size 0..
func (f Function) Inputs() []Parameter {
	return f.inputs
}

// Outputs returns a list of output parameters returned from this function, or nil on invalid object.
// If there are no output parameters, this returns a slice of size 0..
func (f Function) Outputs() []Parameter {
	return f.outputs
}
