// Package pkg provides a generic convenience wrapper around various golang package libraries.
package pkg

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
)

var ErrInvalidPkg = fmt.Errorf("invalid package")

// Package is the main type for this package. It holds details about the package's structure.
type Package struct {
	// Name of this package.
	name string

	// Import path for this package.
	importPath string

	// General package overview comments/documentation.
	comments string

	// List of source files for this package. This includes both the source files for this system's
	// build and those ignored for this system's build.
	files []string

	// List of test files for this package. This includes both the test files within this package
	// and the test files for any other external test package in this package's directory.
	testFiles []string

	// List of imports in the source files (no test files) for this package.
	imports []string

	// List of imports in the test files (no source files) for this package. This includes test
	// files in the package and test files not in the package but in the package's directory.
	testImports []string

	// List of groups of one or more exported constants in this package.
	constantBlocks []ConstantBlock

	// List of groups of one or more exported variables in this package.
	variableBlocks []VariableBlock

	// List of exported functions for this package. This includes only exported functions from the
	// source files, not from the test files.
	functions []Function

	// List of exported types for this package. This includes only exported types from the source
	// files, not from the test files.
	types []Type
}

// New parses the package at importPath and creates a new Package object with its information.
func New(importPath string) (Package, error) {
	// Generate the go/build Package for the import path.
	buildPackage, err := build.Import(importPath, "", 0)
	if err != nil {
		return Package{}, fmt.Errorf("invalid package in %s: %w", importPath, err)
	}

	// Generate all the go/ast Package's for the import path.
	fset := token.NewFileSet()
	astPackages, err := parser.ParseDir(fset, buildPackage.Dir, nil, parser.ParseComments)
	if err != nil {
		return Package{}, fmt.Errorf("invalid package in %s: %w", importPath, err)
	}

	// Get the go/ast Package for the package named by the import path.
	astPackage, ok := astPackages[buildPackage.Name]
	if !ok {
		return Package{}, fmt.Errorf("package not found in %s", importPath)
	}

	// Generate the go/doc Package for the package named by the import path. We first have to
	// flatten out the map of ast files and then use that list to parse the individual files. We
	// want to gather all files for the package and not filter out any based on build system.
	astFiles := make([]*ast.File, 0, len(astPackage.Files))
	for _, v := range astPackage.Files {
		astFiles = append(astFiles, v)
	}
	docPackage, err := doc.NewFromFiles(fset, astFiles, importPath)
	if err != nil {
		return Package{}, fmt.Errorf("invalid package in %s: %w", importPath, err)
	}
	if docPackage == nil {
		return Package{}, ErrInvalidPkg
	}

	// Put everything together into our Package type.
	return newPackage(astPackage, buildPackage, docPackage)
}

// newPackage puts together the internal structure for a Package object.
func newPackage(astPackage *ast.Package, buildPackage *build.Package, docPackage *doc.Package) (Package, error) {
	// Begin with structuring up our object.
	p := Package{
		name:       docPackage.Name,
		importPath: docPackage.ImportPath,
		comments:   docPackage.Doc,
	}

	// Put together the list of source files, both for this system's build and those ignored for
	// this system's build.
	for _, s := range [][]string{buildPackage.GoFiles, buildPackage.IgnoredGoFiles} {
		p.files = append(p.files, s...)
	}
	sort.Strings(p.files)

	// Put together the list of test files, both for this package and any other external test
	// package in this package's directory.
	for _, s := range [][]string{buildPackage.TestGoFiles, buildPackage.XTestGoFiles} {
		p.testFiles = append(p.testFiles, s...)
	}
	sort.Strings(p.testFiles)

	// Copy the list of imports from the source files.
	p.imports = append([]string{}, buildPackage.Imports...)

	// Put together the list of imports from both the internally and externally packaged test files.
	// Make sure the list is sorted and free of duplicates.
	m := make(map[string]bool)
	for _, ss := range [][]string{buildPackage.TestImports, buildPackage.XTestImports} {
		for _, s := range ss {
			m[s] = true
		}
	}
	p.testImports = make([]string, 0, len(m))
	for v := range m {
		p.testImports = append(p.testImports, v)
	}
	sort.Strings(p.testImports)

	// Put together all the source files so we can generate documentation with embedded code.
	r, err := buildSource(buildPackage)
	if err != nil {
		return Package{}, fmt.Errorf("error reading source files: %w", err)
	}

	// Extract the blocks of exported constants for this package, both for standard types (go/doc's
	// Consts) and for custom types (go/doc's Type's Consts).
	for _, cb := range docPackage.Consts {
		p.constantBlocks = append(p.constantBlocks, newConstantBlock(cb, nil, r))
	}
	for _, t := range docPackage.Types {
		for _, cb := range t.Consts {
			p.constantBlocks = append(p.constantBlocks, newConstantBlock(cb, t, r))
		}
	}

	// Extract the blocks of exported variables for this package, both for standard types (go/doc's
	// Vars) and for custom types (go/doc's Type's Vars).
	for _, vb := range docPackage.Vars {
		p.variableBlocks = append(p.variableBlocks, newVariableBlock(vb, nil, r))
	}
	for _, t := range docPackage.Types {
		for _, vb := range t.Vars {
			p.variableBlocks = append(p.variableBlocks, newVariableBlock(vb, t, r))
		}
	}

	// Extract the exported functions for this package.
	p.functions = make([]Function, len(docPackage.Funcs))
	for i, f := range docPackage.Funcs {
		p.functions[i] = newFunction(f, r)
	}

	// Extract the exported types for this package.
	p.types = make([]Type, len(docPackage.Types))
	for i, t := range docPackage.Types {
		p.types[i] = newType(t, r)
	}

	return p, nil
}

// buildSource concatenates the source files for the package into a bytes.Reader.
func buildSource(bp *build.Package) (*bytes.Reader, error) {
	if bp == nil {
		return nil, fmt.Errorf("missing internal build package")
	}

	// First, we need to read in all the source files. When the internal go/* libraries generate the
	// AST and documentation for the package, they read in all source and test files in alphabetical
	// order, joined with newlines. We must use the same approach to make sure that the position
	// indexes line up later.
	files := make([]string, 0)
	for _, s := range [][]string{
		// TODO: what other files need to be added here?
		bp.GoFiles,
		bp.CgoFiles,
		bp.IgnoredGoFiles,
		bp.TestGoFiles,
		bp.XTestGoFiles,
	} {
		files = append(files, s...)
	}
	sort.Strings(files)

	bufs := make([][]byte, len(files))
	for i, f := range files {
		f = filepath.Join(bp.Dir, f)
		data, err := os.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %w", err)
		}
		bufs[i] = data
	}

	return bytes.NewReader(bytes.Join(bufs, []byte{'\n'})), nil
}

// Name returns the package's name.
func (p Package) Name() string {
	return p.name
}

// ImportPath returns the package's import path.
func (p Package) ImportPath() string {
	return p.importPath
}

// Comments returns the general package overview documentation with pkg's formatting applied.
func (p Package) Comments(width int) string {
	return formatComments(p.comments, width)
}

// Files returns a list of source files in the package. The file paths are relative to the package's
// directory, not absolute on the filesystem. Test files (*_test.go) are not included in the list.
// To get a list of test files in the package, see Package's TestFiles. Note: This returns all
// source files in the package's directory and does not limit the files based on what is actually
// used when building for the current system.
func (p Package) Files() []string {
	return append([]string{}, p.files...)
}

// TestFiles returns a list of test files in the package's directory. This includes test files both
// within the package (e.g. mypkg_test.go in package mypkg) and outside of the package but within
// the package's directory (e.g. other_test.go in package mypkg_test). The file paths are relative
// to the package's directory, not absolute on the filesystem. Source files are not included in the
// list. To get a list of source files in the package, see Package's Files.
func (p Package) TestFiles() []string {
	return append([]string{}, p.testFiles...)
}

// Imports returns a list of imports in the package. The list includes only imports from the source
// files, not the test files. To get a list of imports from the test files, see Package's TestImports.
func (p Package) Imports() []string {
	return append([]string{}, p.imports...)
}

// TestImports returns a list of imports from test files both within the package (e.g. mypkg_test.go
// in package mypkg) and outside of the package but within the package's directory (e.g.
// other_test.go in package mypkg_test). The list includes only imports from the test files, not the
// source files. To get a list of imports from the source files, see Package's Imports.
func (p Package) TestImports() []string {
	return append([]string{}, p.testImports...)
}

// ConstantBlocks returns a list of blocks of exported constants in the package. This includes both
// blocks of a standard type (like int or string) and blocks of a custom type (like io.Reader or
// *http.Client). In the latter case, ConstantBlock's Type method can be used to determine the
// block's general type. The list includes only blocks of exported constants from the source files,
// not the test files.
func (p Package) ConstantBlocks() []ConstantBlock {
	return append([]ConstantBlock{}, p.constantBlocks...)
}

// VariableBlocks returns a list of blocks of exported variables in the package. This includes both
// blocks of a standard type (like int or string) and blocks of a custom type (like io.Reader or
// *http.Client). In the latter case, VariableBlock's Type method can be used to determine the
// block's general type. The list includes only blocks of exported variables from the source files,
// not the test files.
func (p Package) VariableBlocks() []VariableBlock {
	return append([]VariableBlock{}, p.variableBlocks...)
}

// Functions returns a list of exported functions in the package. The list includes exported
// functions from source files for the package only, not from test files (internal or external).
func (p Package) Functions() []Function {
	return append([]Function{}, p.functions...)
}

// Types returns a list of exported types in the package.
func (p Package) Types() []Type {
	return append([]Type{}, p.types...)
}
