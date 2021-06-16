// pkg_test is the testing suite for the pkg package. Testing is done by parsing various packages
// and checking that the correct functions, types, methods, etc. are properly read. Each of the
// test packages has been chosen for a specific characteristic, like only having functions or only
// having types without methods. These are the test packages used:
//   errors - Has only functions.
//   fmt - Has only functions and interfaces.
//   hash - Has only types without methods.
//   archive/tar - Has only types with methods, and some global variables/errors.
//   unicode - Has lots of constants and global variables.
//   net/rpc - Has everything.
// TODO: need to also test packages with these features:
//   packages that aren't in the standard library
//   packages with sub-packages
package pkg_test

import (
	"testing"

	"github.com/snhilde/pkg"
)

var (
	testPackages = []string{
		"errors",
		"fmt",
		"hash",
		"archive/tar",
		"unicode",
		"net/rpc",
	}
)
