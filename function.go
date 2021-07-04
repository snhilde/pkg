// This file contains the information and logic for the Function type.
package pkg

import (
	"go/doc"
	"go/token"
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
func newFunction(f *doc.Func, fset *token.FileSet) Function {
	if f == nil {
		return Function{}
	}

	// Extract the parameters.
	in := newParameters(f.Decl.Type.Params, fset)
	out := newParameters(f.Decl.Type.Results, fset)

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
	return append([]Parameter{}, f.inputs...)
}

// Outputs returns a list of output parameters returned from this function, or nil on invalid object.
// If there are no output parameters, this returns a slice of size 0..
func (f Function) Outputs() []Parameter {
	return append([]Parameter{}, f.outputs...)
}
