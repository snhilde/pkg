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
			t.Errorf("%s: incorrect package name", testPackage)
			t.Log("\twant:", want)
			t.Log("\thave:", have)
		}
	}
}

// TestFiles checks that pkg correctly reports the correct source files in each test package.
func TestFiles(t *testing.T) {
	// These are the files in each package. We're going to hard-code these values so that we can
	// achieve repeatable accuracy.
	fileMap := map[string][]string{
		"errors": {
			"errors.go",
			"wrap.go",
		},
		"fmt": {
			"doc.go",
			"errors.go",
			"format.go",
			"print.go",
			"scan.go",
		},
		"hash": {
			"hash.go",
		},
		"archive/tar": {
			"common.go",
			"format.go",
			"reader.go",
			"stat_actime1.go",
			"stat_actime2.go",
			"stat_unix.go",
			"strconv.go",
			"writer.go",
		},
		"unicode": {
			"casetables.go",
			"digit.go",
			"graphic.go",
			"letter.go",
			"tables.go",
		},
		"net/rpc": {
			"client.go",
			"debug.go",
			"server.go",
		},
	}

	for testPackage, files := range fileMap {
		p, _ := pkg.New(testPackage)

		want := files
		have := p.Files()

		// First, let's make sure that we have the correct number of files.
		if len(want) != len(have) {
			t.Errorf("%s: incorrect file list", testPackage)
			t.Log("\twant:", want)
			t.Log("\thave:", have)
		}

		// Then, let's check that each file is present in the returned list.
		for _, w := range want {
			found := false
			for _, h := range have {
				if w == h {
					// If we've already found this file, then something is wrong.
					if found {
						t.Errorf("%s: duplicate file in list: %s", testPackage, w)
						return
					}
					found = true
				}
			}
			if !found {
				t.Errorf("%s: missing file: %s", testPackage, w)
			}
		}
	}
}
