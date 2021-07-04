// This file contains the information and logic for the ConstantBlock and Constant types.
package pkg

import (
	"go/doc"
	"go/token"
)

// ConstantBlock holds information about a block of one or more (grouped) exported constants in a
// package.
type ConstantBlock struct {
	// General type of this block of constants.
	typeName string

	// Comments for this block of constants.
	comments string

	// Original declaration in source for this block of constants.
	source string

	// List of constants within this block.
	constants []Constant
}

// Constant holds information about a single exported constant within a block.
type Constant struct {
	// Name of this constant.
	name string
}

// newConstantBlock builds a new ConstantBlock object based on go/doc's Value.
func newConstantBlock(v *doc.Value, t *doc.Type, fset *token.FileSet) ConstantBlock {
	if v == nil {
		return ConstantBlock{}
	}

	// If this is a block of constants linked to a type, store the name of the type as well.
	typeName := ""
	if t != nil {
		typeName = t.Name
	}

	// Read out the source declaration (exported constants only).
	source := extractSource(v.Decl, fset)

	// Build the list of individual constants.
	constants := make([]Constant, len(v.Names))
	for i, n := range v.Names {
		constants[i] = Constant{name: n}
	}

	return ConstantBlock{
		typeName:  typeName,
		comments:  v.Doc,
		source:    source,
		constants: constants,
	}
}

// Type returns the name of the general, non-built-in type for this block of constants, or "" if this block
// generally does not represent a non-built-in type.
func (cb ConstantBlock) Type() string {
	return cb.typeName
}

// Comments returns the documentation for this block of constants with pkg's formatting applied.
func (cb ConstantBlock) Comments(width int) string {
	return formatComments(cb.comments, width)
}

// Source returns the source declaration for this block of constants.
func (cb ConstantBlock) Source() string {
	return cb.source
}

// Constants returns a list of constants in this block of constants.
func (cb ConstantBlock) Constants() []Constant {
	return append([]Constant{}, cb.constants...)
}

// Name returns the constant's name.
func (c Constant) Name() string {
	return c.name
}
