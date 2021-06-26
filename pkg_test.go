// pkg_test is the testing suite for the pkg package. Testing is done by parsing various packages
// and checking that the correct functions, types, methods, etc. are properly read. Each of the
// test packages has been chosen for a specific characteristic, like only having functions or only
// having types without methods. These are the test packages used:
//   errors - Has only functions.
//   fmt - Has only functions and interfaces. Has backticks in one doc string.
//   hash - Has only types without methods. Makes sure sub-packages are not included. Has exported
//   and unexported variables in the same block of variables.
//   archive/tar - Has only types with methods, and some global variables/errors.
//   unicode - Has lots of constants and global variables and no imported packages. Has unexported
//   constants.
//   net/rpc - Has everything, including a sub-package (net/rpc/jsonrpc), indented package overview
//             comments, and global constants, errors, and variables.
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

		// Check that pkg found the correct blocks of exported constants in the source files.
		if err := cmpConstantBlockLists(testPkg.constantBlocks, p.ConstantBlocks()); err != nil {
			t.Errorf("%s: constant blocks: %s", testPkg.name, err.Error())

			continue
		}

		// Check that pkg found the correct blocks of exported variables and errors in the source files.
		if err := cmpVariableBlockLists(testPkg.variableBlocks, p.VariableBlocks()); err != nil {
			t.Errorf("%s: variable blocks: %s", testPkg.name, err.Error())

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
		return fmt.Errorf("incorrect number of items (want %v, have %v)", len(want), len(have))
	}

	for i, w := range want {
		h := have[i]
		if w != h {
			return fmt.Errorf("mismatch (want %s, have %s)", w, h)
		}
	}

	return nil
}

// cmpConstantBlockLists checks that two lists of blocks of constants have the exact same elements.
func cmpConstantBlockLists(wantConstantBlocks []testConstantBlock, haveConstantBlocks []pkg.ConstantBlock) error {
	if len(wantConstantBlocks) != len(haveConstantBlocks) {
		return fmt.Errorf("incorrect number of constant blocks (want %v, have %v)", len(wantConstantBlocks), len(haveConstantBlocks))
	}

	for i, wantConstantBlock := range wantConstantBlocks {
		haveConstantBlock := haveConstantBlocks[i]

		// Check that the block's type is correct.
		if wantConstantBlock.typeName != haveConstantBlock.Type() {
			return fmt.Errorf("block %v: type mismatch (want %s, have %s)", i, wantConstantBlock.typeName, haveConstantBlock.Type())
		}

		// Check that the block's comments are correct.
		if wantConstantBlock.comments != haveConstantBlock.Comments(99999) {
			return fmt.Errorf("block %v: comments mismatch", i)
		}

		// Check that the block's source is correct.
		if wantConstantBlock.source != haveConstantBlock.Source() {
			return fmt.Errorf("block %v: source mismatch", i)
		}

		// Check that all constants within this block are correct.
		wantConstants := wantConstantBlock.constants
		haveConstants := haveConstantBlock.Constants()
		if len(wantConstants) != len(haveConstants) {
			return fmt.Errorf("block %v: incorrect number of constants (want %v, have %v)", i, len(wantConstants), len(haveConstants))
		}
		for i, wantConstant := range wantConstants {
			haveConstant := haveConstants[i]

			// Check that the constant's name is correct.
			if wantConstant.name != haveConstant.Name() {
				return fmt.Errorf("block %v: constant name mismatch (want %s, have %s)", i, wantConstant.name, haveConstant.Name())
			}
		}
	}

	return nil
}

// cmpVariableBlockLists checks that two lists of blocks of variables and errors have the exact same elements.
func cmpVariableBlockLists(wantVariableBlocks []testVariableBlock, haveVariableBlocks []pkg.VariableBlock) error {
	if len(wantVariableBlocks) != len(haveVariableBlocks) {
		return fmt.Errorf("incorrect number of variable blocks (want %v, have %v)", len(wantVariableBlocks), len(haveVariableBlocks))
	}

	for i, wantVariableBlock := range wantVariableBlocks {
		haveVariableBlock := haveVariableBlocks[i]

		// Check that the block's type is correct.
		if wantVariableBlock.typeName != haveVariableBlock.Type() {
			return fmt.Errorf("block %v: type mismatch (want %s, have %s)", i, wantVariableBlock.typeName, haveVariableBlock.Type())
		}

		// Check that the block's comments are correct.
		if wantVariableBlock.comments != haveVariableBlock.Comments(99999) {
			return fmt.Errorf("block %v: comments mismatch", i)
		}

		// Check that the block's source is correct.
		if wantVariableBlock.source != haveVariableBlock.Source() {
			return fmt.Errorf("block %v: source mismatch", i)
		}

		// Check that all variables within this block are correct.
		wantVariables := wantVariableBlock.variables
		haveVariables := haveVariableBlock.Variables()
		if len(wantVariables) != len(haveVariables) {
			return fmt.Errorf("block %v: incorrect number of variables (want %v, have %v)", i, len(wantVariables), len(haveVariables))
		}
		for i, wantVariable := range wantVariables {
			haveVariable := haveVariables[i]

			// Check that the variable's name is correct.
			if wantVariable.name != haveVariable.Name() {
				return fmt.Errorf("block %v: variable name mismatch (want %s, have %s)", i, wantVariable.name, haveVariable.Name())
			}
		}

		// Check that all errors within this block are correct.
		wantErrors := wantVariableBlock.errors
		haveErrors := haveVariableBlock.Errors()
		if len(wantErrors) != len(haveErrors) {
			return fmt.Errorf("block %v: incorrect number of errors (want %v, have %v)", i, len(wantErrors), len(haveErrors))
		}
		for i, wantError := range wantErrors {
			haveError := haveErrors[i]

			// Check that the error's name is correct.
			if wantError.name != haveError.Name() {
				return fmt.Errorf("block %v: error name mismatch (want %s, have %s)", i, wantError.name, haveError.Name())
			}
		}
	}

	return nil
}

// cmpFunctionLists checks that two lists of functions have the exact same elements.
func cmpFunctionLists(wantFuncs []testFunction, haveFuncs []pkg.Function) error {
	if len(wantFuncs) != len(haveFuncs) {
		return fmt.Errorf("incorrect number of functions (want %v, have %v)", len(wantFuncs), len(haveFuncs))
	}

	for i, wantFunc := range wantFuncs {
		haveFunc := haveFuncs[i]

		// Check that the function's name is correct.
		if wantFunc.name != haveFunc.Name() {
			return fmt.Errorf("name mismatch (want %s, have %s)", wantFunc.name, haveFunc.Name())
		}

		// Check that the function's comments are correct.
		if wantFunc.comments != haveFunc.Comments(99999) {
			return fmt.Errorf("%s: comments mismatch", haveFunc.Name())
		}

		// Check that the input parameters are correct.
		wantInputs := wantFunc.inputs
		haveInputs := haveFunc.Inputs()
		if len(wantInputs) != len(haveInputs) {
			return fmt.Errorf("%s: incorrect number of inputs (want %v, have %v)",
				haveFunc.Name(), len(wantInputs), len(haveInputs))
		}
		for i, wantInput := range wantInputs {
			haveInput := haveInputs[i]

			// Check that the input's name is correct.
			if wantInput.name != haveInput.Name() {
				return fmt.Errorf("%s: input name mismatch (want %s, have %s)",
					haveFunc.Name(), wantInput.name, haveInput.Name())
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
			return fmt.Errorf("%s: incorrect number of outputs (want %v, have %v)",
				haveFunc.Name(), len(wantOutputs), len(haveOutputs))
		}
		for i, wantOutput := range wantOutputs {
			haveOutput := haveOutputs[i]

			// Check that the output's name is correct.
			if wantOutput.name != haveOutput.Name() {
				return fmt.Errorf("%s: s: output name mismatch (want %s, have %s)",
					haveFunc.Name(), wantOutput.name, haveOutput.Name())
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
		return fmt.Errorf("incorrect number of types (want %v, have %v)", len(wantTypes), len(haveTypes))
	}

	for i, wantType := range wantTypes {
		haveType := haveTypes[i]

		// Check that the type's name is correct.
		if wantType.name != haveType.Name() {
			return fmt.Errorf("name mismatch (want %s, have %s)", wantType.name, haveType.Name())
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
		return fmt.Errorf("incorrect number of methods (want %v, have %v)", len(wantMethods), len(haveMethods))
	}

	for i, wantMethod := range wantMethods {
		haveMethod := haveMethods[i]

		// Check that the method's name is correct.
		if wantMethod.name != haveMethod.Name() {
			return fmt.Errorf("name mismatch (want %s, have %s)", wantMethod.name, haveMethod.Name())
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
			return fmt.Errorf("%s: incorrect number of inputs (want %v, have %v)",
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
			return fmt.Errorf("%s: incorrect number of outputs (want %v, have %v)",
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
