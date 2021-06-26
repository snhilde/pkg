// This file contains common functions for this package.
package pkg

import (
	"bytes"
	"go/doc"
	"go/token"
	"io"
	"regexp"
	"strings"
)

// exportedRe is used to find exported constant/variable names.
var exportedRe = regexp.MustCompilePOSIX(`^((const )?|(var )?)[[:upper:]]`)

// extractSource extracts the source text in r between the start and end positions.
func extractSource(r io.ReaderAt, start, end token.Pos) *bytes.Reader {
	if start < 0 || start >= end {
		return nil
	}

	b := make([]byte, end-start)
	if n, err := r.ReadAt(b, int64(start)); n != int(end-start) || err != nil {
		return nil
	}

	return bytes.NewReader(b)
}

// extractGlobalsSource extracts the source declaration for a block of constants/variables.
func extractGlobalsSource(v *doc.Value, r *bytes.Reader, isConst bool) []byte {
	if v == nil || r == nil {
		return nil
	}

	start, end := v.Decl.Pos()-1, v.Decl.End()-1 // -1 to index properly
	decl := extractSource(r, start, end)
	source, err := io.ReadAll(decl)
	if err != nil {
		return nil
	}

	// If there are multiple lines and multiple globals in the source declaration, then it's
	// possible that there is a mix of exported and unexported globals. We need to remove all
	// unexported globals.
	s := strings.Trim(string(source), "\n")
	lines := strings.Split(s, "\n")
	if len(lines) <= 1 {
		// If there's only line in the source, then it has to be an exported global.
		return source
	}

	prefix := ""
	if isConst {
		prefix = "const"
	} else {
		prefix = "var"
	}
	if first := strings.TrimSpace(lines[0]); first != prefix+"(" && first != prefix+" (" {
		// If the line doesn't start with "const(" or "const (" for constants or "var(" or "var ("
		// for variables, then the keyword and global are on the same line, meaning there is only
		// one global in this block and it's exported.
		return source
	}

	// If we're here, the we have multiple lines and multiple globals. To properly restrict the
	// source to only exported globals and their comments, we need to break the source up into its
	// constituent groups, which are separated by blank lines, and then read each group
	// individually.

	// These are the groups that we want to keep in the final source declaration.
	keepGroups := make([]string, 0)

	// Stash the "const (" or "var (" and ")" lines, and remove them from the overall source for
	// easier parsing.
	firstLine := lines[0]
	lastLine := lines[len(lines)-1]
	s = strings.Join(lines[1:len(lines)-1], "\n")

	for _, group := range strings.Split(s, "\n\n") {
		group = extractGroupSource(group)
		if group != "" {
			keepGroups = append(keepGroups, group)
		}
	}

	// Restructure the source with the groups designated for keeping.
	s = strings.Join(keepGroups, "\n\n")

	// Add our first and last lines back in.
	s = strings.Join([]string{firstLine, s, lastLine}, "\n")

	return []byte(s)
}

// extractGroupSource extracts the source declaration for a single group in a block of
// constants/variables. A group is defined as a continuous set of non-blank lines in the source that
// is separated from other sets by blank lines. If the group does not have any exported globals,
// then this returns "".
func extractGroupSource(group string) string {
	// These are the lines that we want to keep for this group.
	keepLines := make([]string, 0)

	// This is the running count of the number of levels of multi-line comment blocks that we're
	// currently in. Every time we find the start of a multi-line comment block (i.e. "/*"), the
	// counter gets incremented, and the opposite for the end of a block ("*/").
	commentBlocks := 0

	// This is the running count of the number of blocks of literals (lines within brackets)
	// we're currently in. Every time we find the start of a block (i.e. "{" at the end of the
	// line), this gets incremented, and the opposite for the end of a block ("}"). If the
	// main/overall global is exported, then saveBlock is true.
	literalBlocks := 0
	saveBlock := false

	// Number of exported globals found in this group.
	numExported := 0

	// Number of unexported globals found in this group.
	numUnexported := 0

	for _, l := range strings.Split(group, "\n") {
		line := strings.TrimSpace(l)

		// Keep all single-line comments.
		if strings.HasPrefix(line, "//") {
			keepLines = append(keepLines, l)

			continue
		}

		if strings.HasPrefix(line, "/*") {
			// This is the start of a (possibly) multi-line comment block.
			commentBlocks++
		}

		// If we're inside a comment block (beginning, middle, or end), then we need to save the
		// line and also check if this is the end of the block.
		if commentBlocks > 0 {
			keepLines = append(keepLines, l)
			if strings.HasSuffix(line, "*/") {
				commentBlocks--
			}

			continue
		}

		if strings.HasSuffix(line, "{") {
			// This is the start of a literal block. If we're not already in a wider block, then
			// let's see if this is an exported or unexported literal.
			if literalBlocks == 0 {
				if exportedRe.MatchString(line) {
					saveBlock = true
					numExported++
				} else {
					saveBlock = false
					numUnexported++
				}
			}
			literalBlocks++
		}

		// If we're inside a literal block (or at its boundary markers), then we need to keep or
		// discard the line based on the overall block's status. We also need to check if this
		// is the last line in the block.
		if literalBlocks > 0 {
			if saveBlock {
				keepLines = append(keepLines, l)
			}

			if strings.HasSuffix(line, "}") {
				// This is the last line in the block.
				literalBlocks--
			}

			continue
		}

		// Finally, check whether or not this line marks an exported global.
		if exportedRe.MatchString(line) {
			numExported++
			keepLines = append(keepLines, l)
		} else {
			numUnexported++
		}
	}

	// We now have all the lines we want to keep for this group. If there are any exported globals
	// in this group, then we want to keep these lines in the overall source.
	if numExported > 0 {
		return strings.Join(keepLines, "\n")
	}

	// We don't have any exported globals in this group. If we also don't have any unexported globals,
	// then maybe this is just an informational group or something. We should probably keep it.
	if numUnexported == 0 {
		return strings.Join(keepLines, "\n")
	}

	// We don't want to keep this group.
	return ""
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
