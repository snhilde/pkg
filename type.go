package pkg

import (
	"bytes"
	"go/doc"
	"io"
)

// Type holds information about an exported type in a package.
type Type struct {
	// Name of this type.
	name string

	// Comments for this type.
	comments string

	// Name of this type's underlying type.
	typeName string

	// Original declaration in source for this type.
	source *bytes.Reader

	// Functions in the package that primarily return this type.
	functions []Function

	// Methods for this type.
	methods []Method
}

// newType builds a new Type object based on go/doc's Type.
func newType(t *doc.Type, r *bytes.Reader) Type {
	if t == nil || r == nil {
		return Type{}
	}
	// Read out the source declaration.
	start, end := t.Decl.Pos()-1, t.Decl.End()-1 // -1 to index properly
	source := extractSource(r, start, end)

	// Extract the underlying type.
	typeName := extractType(t, r)

	// Make a list of functions for this type.
	functions := make([]Function, len(t.Funcs))
	for i, f := range t.Funcs {
		functions[i] = newFunction(f, r)
	}

	// Make a list of methods for this type.
	methods := make([]Method, len(t.Methods))
	for i, m := range t.Methods {
		methods[i] = newMethod(m, r)
	}

	return Type{
		name:      t.Name,
		comments:  t.Doc,
		typeName:  string(typeName),
		source:    source,
		functions: functions,
		methods:   methods,
	}
}

// extractType extracts the underlying type name for this type.
func extractType(t *doc.Type, r *bytes.Reader) []byte {
	// Read out the source declaration.
	start, end := t.Decl.Pos()-1, t.Decl.End()-1 // -1 to index properly
	source := extractSource(r, start, end)
	b, err := io.ReadAll(source)
	if err != nil {
		return nil
	}

	// Remove the type keyword and name of this type from the beginning.
	b = bytes.TrimPrefix(b, []byte("type "+t.Name+" "))

	// Remove any other lines in this declaration.
	if i := bytes.IndexByte(b, '\n'); i > 0 {
		b = b[:i]
	}

	// If we have an opening curly bracket still, we need to remove that.
	b = bytes.TrimSpace(b)
	if bytes.HasSuffix(b, []byte{'{'}) {
		b = b[:len(b)-2]
	}

	return bytes.TrimSpace(b)
}

// Name returns the type's name.
func (t Type) Name() string {
	return t.name
}

// Comments returns the documentation for this type with pkg's formatting applied.
func (t Type) Comments(width int) string {
	return formatComments(t.comments, width)
}

// Type returns the name of the type's underlying type, like "struct" or "map[string]chan int".
func (t Type) Type() string {
	return t.typeName
}

// Functions returns a list of functions that primarily return this type.
func (t Type) Functions() []Function {
	return append([]Function{}, t.functions...)
}

// Methods returns a list of methods for this type.
func (t Type) Methods() []Method {
	return append([]Method{}, t.methods...)
}

func (t Type) Exports() []interface{} {
	return append([]interface{}{}, nil)
}
