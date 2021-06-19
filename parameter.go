// This file contains the logic for extracting and using parameters to/from functions.
package pkg

import (
	"bytes"
	"go/ast"
)

// Parameter represents a parameter to a function, either as input (argument) or output (result).
type Parameter struct {
	// Name of the parameter.
	name string

	// Name of the parameter's type.
	typeName string
}

// extractParameters looks inside ft and extracts all input and output parameters.
func extractParameters(ft *ast.FuncType, r *bytes.Reader) ([]Parameter, []Parameter) {
	if ft == nil || r.Len() == 0 {
		return nil, nil
	}

	// Build the list of inputs.
	in := make([]Parameter, 0)
	if ft.Params != nil {
		for _, p := range ft.Params.List {
			// Read out the type for this group.
			start, end := p.Type.Pos()-1, p.Type.End()-1
			ts := extractSource(r, start, end)
			tb := make([]byte, ts.Len())
			ts.Read(tb)

			// Get the names of all grouped parameters of this type, and add the members to the list.
			for _, name := range p.Names {
				in = append(in, Parameter{
					name:     name.Name,
					typeName: string(tb),
				})
			}
		}
	}

	// Build the list of outputs.
	out := make([]Parameter, 0)
	if ft.Results != nil {
		for _, p := range ft.Results.List {
			// Read out the type for this group.
			start, end := p.Type.Pos()-1, p.Type.End()-1
			ts := extractSource(r, start, end)
			tb := make([]byte, ts.Len())
			ts.Read(tb)

			if p.Names == nil {
				// Unnamed return parameter.
				out = append(out, Parameter{
					name:     "",
					typeName: string(tb),
				})
			} else {
				// Get the names of all grouped parameters of this type and add the members to the list.
				for _, name := range p.Names {
					out = append(out, Parameter{
						name:     name.Name,
						typeName: string(tb),
					})
				}
			}
		}
	}

	return in, out
}

// Name returns the parameter's name. Might be empty for output parameters.
func (p Parameter) Name() string {
	return p.name
}

// Type returns the parameter's type, like "[]byte" or "*os.File".
func (p Parameter) Type() string {
	return p.typeName
}
