package pkg

import (
	"bytes"
	"go/doc"
	"io"
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

// Constant holds information about a single exported constant within a block of (possibly) other
// constants.
type Constant struct {
	// Name of this constant.
	name string
}

// newConstantBlock builds a new ConstantBlock object based on go/doc's Value.
func newConstantBlock(v *doc.Value, t *doc.Type, r *bytes.Reader) ConstantBlock {
	if v == nil || r == nil {
		return ConstantBlock{}
	}

	// If this is a block of constants linked to a type, store the name of the type as well.
	typeName := ""
	if t != nil {
		typeName = t.Name
	}

	// Read out the source declaration.
	start, end := v.Decl.Pos()-1, v.Decl.End()-1 // -1 to index properly
	decl := extractSource(r, start, end)
	source, _ := io.ReadAll(decl)

	constants := make([]Constant, len(v.Names))
	for i, n := range v.Names {
		constants[i] = Constant{name: n}
	}

	return ConstantBlock{
		typeName:  typeName,
		comments:  v.Doc,
		source:    string(source),
		constants: constants,
	}
}

// Type returns the name of the general type for this block of constants.
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
