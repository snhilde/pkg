// pkg_test is the testing suite for the pkg package. Testing is done by parsing various packages
// and checking that the correct functions, types, methods, etc. are properly read. Each of the
// test packages has been chosen for a specific characteristic, like only having functions or only
// having types without methods. These are the test packages used:
//   errors - Has only functions.
//   fmt - Has only functions and interfaces.
//   hash - Has only types without methods.
//   archive/tar - Has only types with methods, and some global variables/errors.
//   unicode - Has lots of constants and global variables and no imported packages.
//   net/rpc - Has everything, including a sub-package (net/rpc/jsonrpc).
// TODO: need to also test packages with these features:
//   packages that aren't in the standard library
//   packages with Cgo files
//   packages with C/C++ files
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
	// These are the source files in each test package. We're going to hard-code these values so
	// that we can achieve repeatable accuracy.
	fileMap := map[string][]string{
		"errors": {
			"errors.go", "wrap.go",
		},
		"fmt": {
			"doc.go", "errors.go", "format.go", "print.go", "scan.go",
		},
		"hash": {
			"hash.go",
		},
		"archive/tar": {
			"common.go", "format.go", "reader.go", "stat_actime1.go", "stat_actime2.go", "stat_unix.go", "strconv.go", "writer.go",
		},
		"unicode": {
			"casetables.go", "digit.go", "graphic.go", "letter.go", "tables.go",
		},
		"net/rpc": {
			"client.go", "debug.go", "server.go",
		},
	}

	checkLists(t, fileMap, pkg.Package.Files)
}

// TestTestFiles checks that pkg correctly reports the correct test files in each test package.
func TestTestFiles(t *testing.T) {
	// These are the test files in each test package. We're going to hard-code these values so that
	// we can achieve repeatable accuracy.
	testFileMap := map[string][]string{
		"errors": {
			"errors_test.go", "example_test.go", "wrap_test.go",
		},
		"fmt": {
			"errors_test.go", "example_test.go", "export_test.go", "fmt_test.go", "gostringer_example_test.go",
			"scan_test.go", "stringer_example_test.go", "stringer_test.go",
		},
		"hash": {
			"example_test.go", "marshal_test.go",
		},
		"archive/tar": {
			"example_test.go", "reader_test.go", "strconv_test.go", "tar_test.go", "writer_test.go",
		},
		"unicode": {
			"digit_test.go", "example_test.go", "graphic_test.go", "letter_test.go", "script_test.go",
		},
		"net/rpc": {
			"client_test.go", "server_test.go",
		},
	}

	checkLists(t, testFileMap, pkg.Package.TestFiles)
}

// TestImports checks that the returned list of imports is correct for each test package.
func TestImports(t *testing.T) {
	// These are the imports used in each test package. We're going to hard-code these values so
	// that we can achieve repeatable accuracy.
	importMap := map[string][]string{
		"errors": {
			"internal/reflectlite",
		},
		"fmt": {
			"errors", "internal/fmtsort", "io", "math", "os", "reflect", "strconv", "sync", "unicode/utf8",
		},
		"hash": {
			"io",
		},
		"archive/tar": {
			"bytes", "errors", "fmt", "io", "io/fs", "math", "os/user", "path", "reflect", "runtime",
			"sort", "strconv", "strings", "sync", "syscall", "time",
		},
		"unicode": {
			// no imports
		},
		"net/rpc": {
			"bufio", "encoding/gob", "errors", "fmt", "go/token", "html/template", "io", "log",
			"net", "net/http", "reflect", "sort", "strings", "sync",
		},
	}

	checkLists(t, importMap, pkg.Package.Imports)
}

// TestTestImports checks that the returned list of test imports is correct for each test package.
func TestTestImports(t *testing.T) {
	// These are the imports used in each test package. We're going to hard-code these values so
	// that we can achieve repeatable accuracy.
	testImportMap := map[string][]string{
		"errors": {
			"errors", "fmt", "io/fs", "os", "reflect", "testing", "time",
		},
		"fmt": {
			"bufio", "bytes", "errors", "fmt", "internal/race", "io", "math", "os", "reflect", "regexp",
			"runtime", "strings", "testing", "testing/iotest", "time", "unicode", "unicode/utf8",
		},
		"hash": {
			"bytes", "crypto/md5", "crypto/sha1", "crypto/sha256", "crypto/sha512", "encoding", "encoding/hex",
			"fmt", "hash", "hash/adler32", "hash/crc32", "hash/crc64", "hash/fnv", "log", "testing",
		},
		"archive/tar": {
			"archive/tar", "bytes", "crypto/md5", "encoding/hex", "errors", "fmt", "internal/testenv", "io",
			"io/fs", "log", "math", "os", "path", "path/filepath", "reflect", "sort", "strconv", "strings",
			"testing", "testing/iotest", "time",
		},
		"unicode": {
			"flag", "fmt", "runtime", "sort", "strings", "testing", "unicode",
		},
		"net/rpc": {
			"errors", "fmt", "io", "log", "net", "net/http/httptest", "reflect", "runtime", "strings",
			"sync", "sync/atomic", "testing", "time",
		},
	}

	checkLists(t, testImportMap, pkg.Package.TestImports)
}

// checkLists checks that the list returned for each of the test packages from method matches the
// expected output in wantMap.
func checkLists(t *testing.T, wantMap map[string][]string, method func(pkg.Package) []string) {
	for testPackage, want := range wantMap {
		p, _ := pkg.New(testPackage)
		have := method(p)

		// First, let's make sure that we have the correct number of items in the list.
		if len(want) != len(have) {
			t.Errorf("%s: incorrect list", testPackage)
			t.Log("\twant:", want)
			t.Log("\thave:", have)
			continue
		}

		// Then, let's check that each wanted item is present in the returned list.
		for _, w := range want {
			found := false
			for _, h := range have {
				if w == h {
					// If we've already found this item, then something is wrong.
					if found {
						t.Errorf("%s: duplicate in list: %s", testPackage, w)
						break
					}
					found = true
				}
			}
			if !found {
				t.Errorf("%s: missing item: %s", testPackage, w)
			}
		}
	}
}
