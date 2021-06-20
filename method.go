// This file contains the information and logic for the Method type.
package pkg

import (
	"bytes"
	"go/doc"
	"strings"
)

// Method holds information about a type's method.
type Method struct {
	// Method's name.
	name string

	// Comments for this method.
	comments string

	// Receiver of this method.
	receiver string

	// Whether or not the receiver of this method is a pointer.
	pointerRcvr bool

	// Input parameters.
	inputs []Parameter

	// Output parameters.
	outputs []Parameter
}

// newMethod builds a new Method object based on go/doc's Func.
func newMethod(f *doc.Func, r *bytes.Reader) Method {
	if f == nil || r == nil {
		return Method{}
	}
	// Read out the source declaration.
	start, end := f.Decl.Type.Pos()-1, f.Decl.Type.End()-1 // -1 to index properly
	decl := extractSource(r, start, end)
	_ = decl // TODO

	// Extract the parameters.
	in, out := extractParameters(f.Decl.Type, r)

	return Method{
		name:        f.Name,
		receiver:    f.Recv,
		pointerRcvr: strings.HasPrefix(f.Recv, "*"),
		inputs:      in,
		outputs:     out,
	}
}

// Name returns the method's name.
func (m Method) Name() string {
	return m.name
}

// Inputs returns a list of input parameters sent to this method, or nil on invalid object. If
// there are no input parameters, this returns a slice of size 0.. The list does not include the
// method receiver.
func (m Method) Inputs() []Parameter {
	return m.inputs
}

// Outputs returns a list of output parameters returned from this method, or nil on invalid object.
// If there are no output parameters, this returns a slice of size 0..
func (m Method) Outputs() []Parameter {
	return m.outputs
}
