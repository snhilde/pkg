// This file contains common functions for this package.
package pkg

import (
	"bytes"
	"go/ast"
	"go/doc"
	"go/token"
	"io"
	"strings"
)

// extractSource extracts the source text in r between the start and end positions.
func extractSource(r io.ReaderAt, start, end token.Pos) string {
	if start < 0 || start >= end {
		return ""
	}

	b := make([]byte, end-start)
	if n, err := r.ReadAt(b, int64(start)); n != int(end-start) || err != nil {
		return ""
	}

	return string(b)
}

// extractExports extracts the source for a block of exported constants or variables. It removes all
// unexported constants/variables and their associated comments from the source and only returns the
// exported constants/variables and their associated comments.
func extractExports(v *doc.Value, r io.ReaderAt) string {
	// In order to extract only the exported constants/variables from this block of source, we're
	// going to descend through the source's AST and pull data out of certain nodes. This is the
	// general approach we're going to take:
	// 1. Split the source into its constituent lines, and cache the line length of each line.
	// 2. Descend through the AST, according to this pattern:
	//     a. The top node is always an *ast.GenDecl. If this node reports that the block is *not*
	//        parenthesized, then we have only constant/variable block, and it must be exported.
	//        We'll return the full source of the block. Otherwise, we'll continue descending.
	//     b. An *ast.ValueSpec represents a single line of one or more declarations, including their
	//        associated comments. When we arrive at one of these nodes, we're first going to check
	//        for the associated comments. If those are present, then we know that we need to start
	//        a new set by inserting a blank line, and then extracting the comments appropriately.
	//        Unfortunately, we don't have a way to read out the entire line of source directly;
	//        there are only accessors for reconstructing the parts individually. Instead, we're
	//        going to figure out where the line begins and ends, and then slice that out.
	// 3. If the block is parenthesized, then we'll add the first and last lines back in.
	if v == nil || v.Decl == nil {
		return ""
	}

	// This is only for constants and variables.
	if v.Decl.Tok != token.CONST && v.Decl.Tok != token.VAR {
		return ""
	}

	// Remove all unexported nodes.
	ast.FilterDecl(v.Decl, ast.IsExported)

	start, end := v.Decl.Pos()-1, v.Decl.End()-1 // -1 to index properly
	source := extractSource(r, start, end)
	offset := start

	parenthesized := false
	sb := new(strings.Builder)
	ast.Inspect(v.Decl, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.GenDecl:
			parenthesized = node.Lparen.IsValid()
			if !parenthesized {
				sb.WriteString(source)
				return false
			}
			return true
		case *ast.ValueSpec:
			if doc := node.Doc; doc != nil {
				if sb.Len() > 0 {
					sb.WriteString("\n")
				}
				for _, comment := range doc.List {
					sb.WriteString("\n\t" + comment.Text)
				}
			}

			start := node.Names[0].Pos()-1 // -1 to index properly
			length := strings.Index(source[start-offset:], "\n")
			end := start + token.Pos(length)
			line := extractSource(r, start, end)
			sb.WriteString("\n\t" + line)
		}
		return false

	})

	s := sb.String()
	if parenthesized {
		lines := strings.Split(source, "\n")
		s = lines[0] + s + "\n" + lines[len(lines)-1]
	}

	return s
}

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
