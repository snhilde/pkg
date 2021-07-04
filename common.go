// This file contains common functions for this package.
package pkg

import (
	"bytes"
	"go/ast"
	"go/doc"
	"go/format"
	"go/token"
	"io"
	"strings"
)

// formatComments returns formatted comment text with pkg's formatting applied.
func formatComments(comments string, width int) string {
	if width <= 0 {
		return comments
	}

	b := new(bytes.Buffer)
	doc.ToText(b, comments, "", "\t", width)

	r, err := io.ReadAll(b)
	if err != nil {
		return ""
	}

	return string(r)
}

// extractSource extracts the source for a node's declaration.
func extractSource(node ast.Node, fset *token.FileSet) string {
	if node == nil || fset == nil {
		return ""
	}

	sb := new(strings.Builder)
	if err := format.Node(sb, fset, node); err != nil {
		return ""
	}

	s := sb.String()

	// Remove multiple blank lines.
	for strings.Contains(s, "\n\n\n") {
		s = strings.ReplaceAll(s, "\n\n\n", "\n\n")
	}

	// If there's a blank line after the start of an identifier block (like can happen when
	// unexported items are removed), remove the blank line.
	s = strings.ReplaceAll(s, "{\n\n", "{\n")
	s = strings.ReplaceAll(s, "(\n\n", "(\n")

	return s
}
