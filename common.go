// This file contains common functions for this package.
package pkg

import (
	"bytes"
	"go/token"
	"io"
)

// extractSource extracts the source text in r between the start and end positions.
func extractSource(r io.ReaderAt, start, end token.Pos) *bytes.Reader {
	if start < 0 || start >= end {
		return nil
	}

	b := make([]byte, end - start)
	if n, err := r.ReadAt(b, int64(start)); n != int(end - start) || err != nil {
		return nil
	}

	return bytes.NewReader(b)
}
