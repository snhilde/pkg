// This file contains the information and logic for the Method type.
package pkg

import (
	"go/doc"
	"go/token"
)

// Method holds information about a type's method.
type Method struct {
	// Method's name.
	name string

	// Comments for this method.
	comments string

	// Receiver of this method.
	receiver Parameter

	// Input parameters.
	inputs []Parameter

	// Output parameters.
	outputs []Parameter
}

// newMethod builds a new Method object based on go/doc's Func.
func newMethod(m *doc.Func, fset *token.FileSet) Method {
	if m == nil {
		return Method{}
	}

	// Extract the receiver.
	var receiver Parameter
	receivers := newParameters(m.Decl.Recv, fset)
	if len(receivers) > 0 {
		receiver = receivers[0]
	}

	// Extract the parameters.
	in := newParameters(m.Decl.Type.Params, fset)
	out := newParameters(m.Decl.Type.Results, fset)

	return Method{
		name:     m.Name,
		comments: m.Doc,
		receiver: receiver,
		inputs:   in,
		outputs:  out,
	}
}

// Name returns the method's name.
func (m Method) Name() string {
	return m.name
}

// Comments returns the documentation for this method with pkg's formatting applied.
func (m Method) Comments(width int) string {
	return formatComments(m.comments, width)
}

// Receiver returns method's receiver.
func (m Method) Receiver() Parameter {
	return m.receiver
}

// PointerReceiver returns whether or not the method's receiver is a pointer to the method's type.
func (m Method) PointerReceiver() bool {
	return m.receiver.Pointer()
}

// Inputs returns a list of input parameters sent to this method, or nil on invalid object. If
// there are no input parameters, this returns a slice of size 0.. The list does not include the
// method's receiver.
func (m Method) Inputs() []Parameter {
	return append([]Parameter{}, m.inputs...)
}

// Outputs returns a list of output parameters returned from this method, or nil on invalid object.
// If there are no output parameters, this returns a slice of size 0..
func (m Method) Outputs() []Parameter {
	return append([]Parameter{}, m.outputs...)
}
