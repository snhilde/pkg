// pkg_test is the testing suite for the pkg package. Testing is done by parsing various packages
// and checking that the correct functions, types, methods, etc. are properly read. Each of the
// test packages has been chosen for a specific characteristic, like only having functions or only
// having types without methods. These are the test packages used:
//   errors - Has only functions.
//   fmt - Has only functions and interfaces. Has backticks in one doc string.
//   hash - Has only types without methods.
//   archive/tar - Has only types with methods, and some global variables/errors.
//   unicode - Has lots of constants and global variables and no imported packages.
//   net/rpc - Has everything, including a sub-package (net/rpc/jsonrpc) and indented package
//             overview comments.
// TODO: need to also test packages with these features:
//   packages that aren't in the standard library
//   packages with Cgo files
//   packages with C/C++ files
//   packages with file naming conventions that might sort differently
// TODO: need invalid object tests
package pkg_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/snhilde/pkg"
)

var testPackages = []testPackage{
	pkgErrors,
	pkgFmt,
	pkgHash,
	pkgArchiveTar,
	pkgUnicode,
	pkgNetRPC,
}

func TestPackages(t *testing.T) {
	t.Parallel()

	for _, testPkg := range testPackages {
		// Check that we can create a new Package with this import path.
		p, err := pkg.New(testPkg.importPath)
		if err != nil {
			t.Error(err)

			continue
		}

		// Check that the package name is correct. The package name is the last member in the import path.
		members := strings.Split(testPkg.name, "/")
		want := members[len(members)-1]
		have := p.Name()
		if want != have {
			t.Errorf("%s: incorrect package name", testPkg.name)
			t.Log("\twant:", want)
			t.Log("\thave:", have)
		}

		// Check that the package's import path is correct.
		want = testPkg.importPath
		have = p.ImportPath()
		if want != have {
			t.Errorf("%s: incorrect package import path", testPkg.name)
			t.Log("\twant:", want)
			t.Log("\thave:", have)
		}

		// Check that the general package overview comments are correct.
		want = testPkg.comments
		have = p.Comments(99999)
		if want != have {
			t.Errorf("%s: incorrect package comments", testPkg.name)
			t.Log("\twant:\n", want)
			t.Log("\thave:\n", have)

			continue
		}

		// Check that pkg found the correct source files in the package's directory and that
		// everything is returned in order with no duplicates.
		if err := cmpStringLists(testPkg.files, p.Files()); err != nil {
			t.Errorf("%s: source files: %s", testPkg.name, err.Error())

			continue
		}

		// Check that pkg found the correct test files in the package's directory and that
		// everything is returned in order with no duplicates.
		if err := cmpStringLists(testPkg.testFiles, p.TestFiles()); err != nil {
			t.Errorf("%s: test files: %s", testPkg.name, err.Error())

			continue
		}

		// Check that pkg found the correct imports in the source files and that everything is
		// returned in order with no duplicates.
		if err := cmpStringLists(testPkg.imports, p.Imports()); err != nil {
			t.Errorf("%s: source imports: %s", testPkg.name, err.Error())

			continue
		}

		// Check that pkg found the correct imports in the test files and that everything is
		// returned in order with no duplicates.
		if err := cmpStringLists(testPkg.testImports, p.TestImports()); err != nil {
			t.Errorf("%s: test imports: %s", testPkg.name, err.Error())

			continue
		}

		// Check that pkg found the correct functions in the source files and that everything is
		// returned in order with no duplicates.
		if err := cmpFunctionLists(testPkg.functions, p.Functions()); err != nil {
			t.Errorf("%s: functions: %s", testPkg.name, err.Error())

			continue
		}

		// Check that pkg found the correct types in the source files and that everything is
		// returned in order with no duplicates.
		if err := cmpTypeLists(testPkg.types, p.Types()); err != nil {
			t.Errorf("%s: types: %s", testPkg.name, err.Error())

			continue
		}
	}
}

// cmpLists checks that two lists of strings have the exact same elements.
func cmpStringLists(want, have []string) error {
	if len(want) != len(have) {
		return fmt.Errorf("incorrent number of items (want %v, have %v)", len(want), len(have))
	}

	for i, w := range want {
		h := have[i]
		if w != h {
			return fmt.Errorf("mismatch (want %s, have %s)", w, h)
		}
	}

	return nil
}

// cmpFunctionLists checks that two lists of functions have the exact same elements.
func cmpFunctionLists(wantFuncs []testFunction, haveFuncs []pkg.Function) error {
	if len(wantFuncs) != len(haveFuncs) {
		return fmt.Errorf("incorrent number of functions (want %v, have %v)", len(wantFuncs), len(haveFuncs))
	}

	for i, wantFunc := range wantFuncs {
		haveFunc := haveFuncs[i]

		// Check that the function's name is correct.
		if wantFunc.name != haveFunc.Name() {
			return fmt.Errorf("%s: name mismatch (want %s, have %s)", haveFunc.Name(), wantFunc.name, haveFunc.Name())
		}

		// Check that the function's comments are correct.
		if wantFunc.comments != haveFunc.Comments(99999) {
			return fmt.Errorf("%s: comments mismatch", haveFunc.Name())
		}

		// Check that the input parameters are correct.
		wantInputs := wantFunc.inputs
		haveInputs := haveFunc.Inputs()
		if len(wantInputs) != len(haveInputs) {
			return fmt.Errorf("%s: incorrent number of inputs (want %v, have %v)",
				haveFunc.Name(), len(wantInputs), len(haveInputs))
		}
		for i, wantInput := range wantInputs {
			haveInput := haveInputs[i]

			// Check that the input's name is correct.
			if wantInput.name != haveInput.Name() {
				return fmt.Errorf("%s: %s: input name mismatch (want %s, have %s)",
					haveFunc.Name(), haveInput.Name(), wantInput.name, haveInput.Name())
			}

			// Check that the input type's name is correct.
			if wantInput.typeName != haveInput.Type() {
				return fmt.Errorf("%s: %s: input type mismatch (want %s, have %s)",
					haveFunc.Name(), haveInput.Name(), wantInput.typeName, haveInput.Type())
			}
		}

		// Check that the output parameters are correct.
		wantOutputs := wantFunc.outputs
		haveOutputs := haveFunc.Outputs()
		if len(wantOutputs) != len(haveOutputs) {
			return fmt.Errorf("%s: incorrent number of outputs (want %v, have %v)",
				haveFunc.Name(), len(wantOutputs), len(haveOutputs))
		}
		for i, wantOutput := range wantOutputs {
			haveOutput := haveOutputs[i]

			// Check that the output's name is correct.
			if wantOutput.name != haveOutput.Name() {
				return fmt.Errorf("%s: %s: output name mismatch (want %s, have %s)",
					haveFunc.Name(), haveOutput.Name(), wantOutput.name, haveOutput.Name())
			}

			// Check that the output type's name is correct.
			if wantOutput.typeName != haveOutput.Type() {
				return fmt.Errorf("%s: %s: output type mismatch (want %s, have %s)",
					haveFunc.Name(), haveOutput.Name(), wantOutput.typeName, haveOutput.Type())
			}
		}
	}

	return nil
}

// cmpTypeLists checks that two lists of types have the exact same elements.
func cmpTypeLists(wantTypes []testType, haveTypes []pkg.Type) error {
	if len(wantTypes) != len(haveTypes) {
		return fmt.Errorf("incorrent number of types (want %v, have %v)", len(wantTypes), len(haveTypes))
	}

	for i, wantType := range wantTypes {
		haveType := haveTypes[i]

		// Check that the type's name is correct.
		if wantType.name != haveType.Name() {
			return fmt.Errorf("%s: name mismatch (want %s, have %s)", haveType.Name(), wantType.name, haveType.Name())
		}

		// Check that the type's comments are correct.
		if wantType.comments != haveType.Comments(99999) {
			return fmt.Errorf("%s: comments mismatch", haveType.Name())
		}

		// Check that the name of the type's underlying type is correct.
		if wantType.typeName != haveType.Type() {
			return fmt.Errorf("%s: type mismatch (want %s, have %s)", haveType.Name(), wantType.typeName, haveType.Type())
		}

		// Check that the type's source is correct.
		// TODO

		// Check that the type's functions are correct.
		if err := cmpFunctionLists(wantType.functions, haveType.Functions()); err != nil {
			return fmt.Errorf("%s: %w", haveType.Name(), err)
		}

		// Check that the type's methods are correct.
		if err := cmpMethodLists(wantType.methods, haveType.Methods()); err != nil {
			return fmt.Errorf("%s: %w", haveType.Name(), err)
		}
	}

	return nil
}

// cmpMethodLists checks that two lists of methods have the exact same elements.
func cmpMethodLists(wantMethods []testMethod, haveMethods []pkg.Method) error {
	if len(wantMethods) != len(haveMethods) {
		return fmt.Errorf("incorrent number of methods (want %v, have %v)", len(wantMethods), len(haveMethods))
	}

	for i, wantMethod := range wantMethods {
		haveMethod := haveMethods[i]

		// Check that the method's name is correct.
		if wantMethod.name != haveMethod.Name() {
			return fmt.Errorf("%s: name mismatch (want %s, have %s)", haveMethod.Name(), wantMethod.name, haveMethod.Name())
		}

		// Check that the method's comments are correct.
		if wantMethod.comments != haveMethod.Comments(99999) {
			return fmt.Errorf("%s: comments mismatch", haveMethod.Name())
		}

		// Check that the method's receiver is correct.
		if wantMethod.receiver != haveMethod.Receiver() {
			return fmt.Errorf("%s: receiver mismatch (want %s, have %s)",
				haveMethod.Name(), wantMethod.receiver, haveMethod.Receiver())
		}

		// Check that the method's receiver's pointer status is correct.
		if wantMethod.pointerRcvr != haveMethod.PointerReceiver() {
			return fmt.Errorf("%s: pointer receiver mismatch (want %v, have %v)",
				haveMethod.Name(), wantMethod.pointerRcvr, haveMethod.PointerReceiver())
		}

		// Check that the input parameters are correct.
		wantInputs := wantMethod.inputs
		haveInputs := haveMethod.Inputs()
		if len(wantInputs) != len(haveInputs) {
			return fmt.Errorf("%s: incorrent number of inputs (want %v, have %v)",
				haveMethod.Name(), len(wantInputs), len(haveInputs))
		}
		for i, wantInput := range wantInputs {
			haveInput := haveInputs[i]

			// Check that the input's name is correct.
			if wantInput.name != haveInput.Name() {
				return fmt.Errorf("%s: input name mismatch (want %s, have %s)",
					haveMethod.Name(), wantInput.name, haveInput.Name())
			}

			// Check that the input type's name is correct.
			if wantInput.typeName != haveInput.Type() {
				return fmt.Errorf("%s: input type mismatch (want %s, have %s)",
					haveMethod.Name(), wantInput.typeName, haveInput.Type())
			}
		}

		// Check that the output parameters are correct.
		wantOutputs := wantMethod.outputs
		haveOutputs := haveMethod.Outputs()
		if len(wantOutputs) != len(haveOutputs) {
			return fmt.Errorf("%s: incorrent number of outputs (want %v, have %v)",
				haveMethod.Name(), len(wantOutputs), len(haveOutputs))
		}
		for i, wantOutput := range wantOutputs {
			haveOutput := haveOutputs[i]

			// Check that the output's name is correct.
			if wantOutput.name != haveOutput.Name() {
				return fmt.Errorf("%s: output name mismatch (want %s, have %s)",
					haveMethod.Name(), wantOutput.name, haveOutput.Name())
			}

			// Check that the output type's name is correct.
			if wantOutput.typeName != haveOutput.Type() {
				return fmt.Errorf("%s: output type mismatch (want %s, have %s)",
					haveMethod.Name(), wantOutput.typeName, haveOutput.Type())
			}
		}
	}

	return nil
}
