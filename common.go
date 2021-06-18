// This file contains common functions for this package.
package pkg

import (
	"bytes"
	"io"
)

// getSource gets the source text in r between the start and end positions, optionally adding a
// prefix to the front of it.
func getSource(r *bytes.Reader, start, end int, prefix []byte) []byte {
	if start < 0 || start >= end || end >= r.Len() {
		return nil
	}

	// Move to the beginning of this segment of code.
	if n, err := r.Seek(int64(start), io.SeekStart); n != int64(start) || err != nil {
		return nil
	}

	// Read in the source.
	b := make([]byte, end - start)
	if n, err := r.Read(b); n != end - start || err != nil {
		return nil
	}

	// If a prefix is present, add that to the beginning.
	return append(prefix, b...)
}
