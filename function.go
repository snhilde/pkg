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

	// Arguments sent in to the function.
	args []Parameter

	// Results returned from the functions.
	rets []Parameter
}

// extractFunctions returns a list of exported functions from source files in p.
func (p Package) extractFunctions(r *bytes.Reader) []Function {
	if !p.IsValid() {
		return nil
	}

	funcs := make([]Function, len(p.docPackage.Funcs))
	for i, v := range p.docPackage.Funcs {
		tp := v.Decl.Type

		// Read out the source declaration.
		decl := extractSource(r, start, end)

		// Extract the parameters.
		in, out := extractParameters(tp)

		start, end := tp.Pos()-1, tp.End()-1 // -1 to index properly
		funcs[i] = Function{
			docFunc: v,
			name:    v.Name,
			decl:    decl,
			args:    in,
			rets:    out,
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
