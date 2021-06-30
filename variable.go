// This file contains the information and logic for the VariableBlock, Variable, and Error types.
package pkg

import (
	"bytes"
	"go/doc"
	"regexp"
)

// errRe is used to find variable names that are exported errors.
var errRe = regexp.MustCompilePOSIX(`^Err[[:upper:]][[:word:]]*$`)

// VariableBlock holds information about a block of one or more (grouped) exported variables in a
// package.
type VariableBlock struct {
	// General type of this block of variables.
	typeName string

	// Comments for this block of variables.
	comments string

	// Original declaration in source for this block of variables.
	source string

	// List of variables within this block.
	variables []Variable

	// List of errors within this block.
	errors []Error
}

// Variable holds information about a single exported variable within a block.
type Variable struct {
	// Name of this variable.
	name string
}

// Error holds information about a single exported error within a block.
type Error struct {
	// Name of this error.
	name string
}

// newVariableBlock builds a new VariableBlock object based on go/doc's Value.
func newVariableBlock(v *doc.Value, t *doc.Type, r *bytes.Reader) VariableBlock {
	if v == nil || r == nil {
		return VariableBlock{}
	}

	// If this is a block of variables linked to a type, store the name of the type as well.
	typeName := ""
	if t != nil {
		typeName = t.Name
	}

	// Read out the source declaration.
	source := extractExports(v, r)

	// Save the variable/error names.
	variables := make([]Variable, len(v.Names))
	errors := make([]Error, 0)
	for i, n := range v.Names {
		variables[i] = Variable{name: n}

		// If this is an error, add it to the list of errors in this block.
		if errRe.MatchString(n) {
			errors = append(errors, Error{name: n})
		}
	}

	return VariableBlock{
		typeName:  typeName,
		comments:  v.Doc,
		source:    source,
		variables: variables,
		errors:    errors,
	}
}

// Type returns the name of the general, non-built-in type for this block of variables, or "" if this block
// generally does not represent a non-built-in type.
func (vb VariableBlock) Type() string {
	return vb.typeName
}

// Comments returns the documentation for this block of variables with pkg's formatting applied.
func (vb VariableBlock) Comments(width int) string {
	return formatComments(vb.comments, width)
}

// Source returns the source declaration for this block of variables.
func (vb VariableBlock) Source() string {
	return vb.source
}

// Variables returns a list of variables in this block of variables. The list includes variables of
// type "error".
func (vb VariableBlock) Variables() []Variable {
	return append([]Variable{}, vb.variables...)
}

// Errors returns a list of only variables of type "error" in this block of variables.
func (vb VariableBlock) Errors() []Error {
	return append([]Error{}, vb.errors...)
}

// Name returns the variable's name.
func (v Variable) Name() string {
	return v.name
}

// Name returns the error's name.
func (e Error) Name() string {
	return e.name
}
