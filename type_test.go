// See pkg_test.go for main testing documentation.
// This file contains the tests for pkg.Test.
package pkg_test

import (
	"testing"

	"github.com/snhilde/pkg"
)

// TODO: need invalid object tests

// TestTypeFunctions checks that the returned list of functions returning each type is correct for
// each of the test packages.
func TestTypeFunctions(t *testing.T) {
	t.Parallel()

	// These are the functions returning each type in each test package. We're going to hard-code
	// these values so that we can achieve repeatable accuracy.
	funcMap := map[string]map[string][]string{
		"errors": {
			// no types
		},
		"fmt": {
			// no functions for any of the types
		},
		"hash": {
			// no functions for any of the types
		},
		"archive/tar": {
			"Header": {
				"FileInfoHeader",
			},
			"Reader": {
				"NewReader",
			},
			"Writer": {
				"NewWriter",
			},
		},
		"unicode": {
			// no functions for any of the types
		},
		"net/rpc": {
			"Client": {
				"Dial", "DialHTTP", "DialHTTPPath", "NewClient", "NewClientWithCodec",
			},
			"Server": {
				"NewServer",
			},
		},
	}

	checkTypeItems(t, funcMap, func(tp pkg.Type) []string {
		funcs := tp.Functions()
		funcNames := make([]string, len(funcs))
		for i, f := range funcs {
			funcNames[i] = f.Name()
		}

		return funcNames
	})
}

// TestTypeMethods checks that the returned list of methods for each type is correct for each of the
// test packages.
func TestTypeMethods(t *testing.T) {
	t.Parallel()

	// These are the methods for each type in each test package. We're going to hard-code these
	// values so that we can achieve repeatable accuracy.
	methodMap := map[string]map[string][]string{
		"errors": {
			// no types
		},
		"fmt": {
			// no methods for any of the types
		},
		"hash": {
			// no methods for any of the types
		},
		"archive/tar": {
			"Format": {
				"String",
			},
			"Header": {
				"FileInfo",
			},
			"Reader": {
				"Next", "Read",
			},
			"Writer": {
				"Close", "Flush", "Write", "WriteHeader",
			},
		},
		"unicode": {
			"SpecialCase": {
				"ToLower", "ToTitle", "ToUpper",
			},
		},
		"net/rpc": {
			"Client": {
				"Call", "Close", "Go",
			},
			"Server": {
				"Accept", "HandleHTTP", "Register", "RegisterName", "ServeCodec", "ServeConn", "ServeHTTP", "ServeRequest",
			},
			"ServerError": {
				"Error",
			},
		},
	}

	checkTypeItems(t, methodMap, func(tp pkg.Type) []string {
		methods := tp.Methods()
		methodNames := make([]string, len(methods))
		for i, m := range methods {
			methodNames[i] = m.Name()
		}

		return methodNames
	})
}

// checkTypeItems checks that the list returned for each type in each of the test packages matches
// the expected output in pkgMap. This is used for getting a list from a list, like checking all the
// methods for all the types in a test package.
func checkTypeItems(t *testing.T, pkgMap map[string]map[string][]string, cb func(pkg.Type) []string) {
	t.Helper()

	for _, testPackage := range testPackages {
		p, _ := pkg.New(testPackage)
		wantTypes, ok := pkgMap[testPackage]
		if !ok {
			t.Errorf("%s: missing from test map", testPackage)

			continue
		}

		// Iterate through each of the types found in this test package.
		for _, haveType := range p.Types() {
			wantItems := wantTypes[haveType.Name()]
			haveItems := cb(haveType)

			// First, let's make sure that we have the correct number of items in the list for
			// this type.
			if len(wantItems) != len(haveItems) {
				t.Errorf("%s: incorrect list", testPackage)
				t.Log("\twant:", wantItems)
				t.Log("\thave:", haveItems)

				break
			}

			// Then, let's check that each type in the list returned the correct items.
			for _, w := range wantItems {
				found := false
				for _, h := range haveItems {
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
}
