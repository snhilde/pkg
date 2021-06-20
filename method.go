// This file contains the information and logic for the Method type.
package pkg

import (
	"bytes"
	"go/doc"
	"io"
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

	// Extract the receiver.
	receiver, pointerRcvr := extractReceiver(f, r)

	// Extract the parameters.
	in, out := extractParameters(f.Decl.Type, r)

	return Method{
		name:        f.Name,
		receiver:    string(receiver),
		pointerRcvr: pointerRcvr,
		inputs:      in,
		outputs:     out,
	}
}

// extractReceiver parses the source of this method and extracts its receiver, also returning
// whether or not the receiver is a pointer.
func extractReceiver(f *doc.Func, r *bytes.Reader) ([]byte, bool) {
	// Read out the source declaration.
	start, end := f.Decl.Type.Pos()-1, f.Decl.Type.End()-1 // -1 to index properly
	source := extractSource(r, start, end)
	b, err := io.ReadAll(source)
	if err != nil {
		return nil, false
	}

	// Remove the func keyword and opening parenthesis from the beginning.
	b = bytes.TrimPrefix(b, []byte("func ("))

	// Remove everything after (and including) the closing parenthesis.
	i := bytes.IndexByte(b, ')')
	b = bytes.TrimSpace(b[:i])

	// Figure out if the receiver is a pointer or not.
	i = bytes.IndexByte(b, ' ')
	isPointer := b[i+1] == '*'

	return b, isPointer
}

// Name returns the method's name.
func (m Method) Name() string {
	return m.name
}

// Receiver returns method's receiver.
func (m Method) Receiver() string {
	return m.receiver
}

// PointerReceiver returns whether or not the method's receiver is a pointer to the method's type.
func (m Method) PointerReceiver() bool {
	return m.pointerRcvr
}

// Inputs returns a list of input parameters sent to this method, or nil on invalid object. If
// there are no input parameters, this returns a slice of size 0.. The list does not include the
// method's receiver.
func (m Method) Inputs() []Parameter {
	return m.inputs
}

// Outputs returns a list of output parameters returned from this method, or nil on invalid object.
// If there are no output parameters, this returns a slice of size 0..
func (m Method) Outputs() []Parameter {
	return m.outputs
}
