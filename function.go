package pkg

import (
	"bytes"
	"go/doc"
)

// Function holds information about an exported function in a package.
type Function struct {
	// Func object from go/doc.
	docFunc *doc.Func

	// Function name.
	name string

	// Original declaration in source for this function.
	decl *bytes.Reader

	// Input parameters.
	input []Parameter

	// Output parameters.
	output []Parameter
}

// extractFunctions returns a list of exported functions from source files in p.
func (p Package) extractFunctions(r *bytes.Reader) []Function {
	if !p.IsValid() {
		return nil
	}

	funcs := make([]Function, len(p.docPackage.Funcs))
	for i, v := range p.docPackage.Funcs {
		// Read out the source declaration.
		start, end := v.Decl.Type.Pos()-1, v.Decl.Type.End()-1 // -1 to index properly
		decl := extractSource(r, start, end)

		// Extract the parameters.
		in, out := extractParameters(v.Decl, r)

		// Put everything together.
		funcs[i] = Function{
			docFunc: v,
			name:    v.Name,
			decl:    decl,
			input:   in,
			output:  out,
		}
	}

	return funcs
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
