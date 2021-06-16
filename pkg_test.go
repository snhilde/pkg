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
	"strings"
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

// TestNew tests creating a new Package for each of the test packages.
func TestNew(t *testing.T) {
	for _, testPackage := range testPackages {
		p, err := pkg.New(testPackage)
		if err != nil {
			t.Error(err)
		} else if p == (pkg.Package{}) {
			t.Errorf("%s: received empty Package object", testPackage)
		}
	}
}

// TestName checks that the name is correct for each test package.
func TestName(t *testing.T) {
	for _, testPackage := range testPackages {
		p, _ := pkg.New(testPackage)

		// The package name is the last member in the import path.
		members := strings.Split(testPackage, "/")
		want := members[len(members)-1]
		have := p.Name()
		if want != have {
			t.Error("incorrect package name")
			t.Log("\twant:", want)
			t.Log("\thave:", have)
		}
	}
}
