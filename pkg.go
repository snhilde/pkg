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
	// Package object from go/ast. This is used to gather the files in the particular package we
	// want and exclude other packages in the same import directory (like external test packages).
	astPackage *ast.Package

	// Package object from go/build. This holds information about the files in the import directory,
	// including source files, test files in the package, and test files not in the package but in
	// the import directory. It does not hold much information about the code structure.
	buildPackage *build.Package

	// Package object from go/doc. This holds the information about the structure of the code.
	docPackage *doc.Package

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
	astPackages, err := parser.ParseDir(fset, buildPackage.Dir, nil, parser.AllErrors)
	if err != nil {
		return Package{}, fmt.Errorf("invalid package in %s: %w", importPath, err)
	}

	// Get the go/ast Package for the package named by the import path.
	astPackage, ok := astPackages[buildPackage.Name]
	if !ok {
		return Package{}, fmt.Errorf("package not found in %s", importPath)
	}

	// Generate the go/doc Package for the package named by the import path. We first have to
	// flatten out the map of ast files and then use that list to parse the individual files.
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
		astPackage:   astPackage,
		buildPackage: buildPackage,
		docPackage:   docPackage,
	}

	// Put together the list of source files, both for this system's build and those ignored for
	// this system's build.
	for _, s := range [][]string{p.buildPackage.GoFiles, p.buildPackage.IgnoredGoFiles} {
		p.files = append(p.files, s...)
	}

	// Put together the list of test files, both for this package and any other external test
	// package in this package's directory.
	for _, s := range [][]string{p.buildPackage.TestGoFiles, p.buildPackage.XTestGoFiles} {
		p.testFiles = append(p.testFiles, s...)
	}

	// Copy the list of imports from the source files.
	p.imports = append([]string{}, p.buildPackage.Imports...)

	// Put together the list of imports from both the internally and externally packaged test files.
	// Make sure the list is sorted and free of duplicates.
	m := make(map[string]bool)
	for _, ss := range [][]string{p.buildPackage.TestImports, p.buildPackage.XTestImports} {
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

	// Extract the exported functions for this package.
	p.functions = make([]Function, len(p.docPackage.Funcs))
	for i, f := range p.docPackage.Funcs {
		p.functions[i] = newFunction(f, r)
	}

	// Extract the exported types for this package.
	p.types = make([]Type, len(p.docPackage.Types))
	for i, t := range p.docPackage.Types {
		p.types[i] = newType(t, r)
	}

	// TODO: extract all other information

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

// IsValid checks whether or not p is a valid Package object.
func (p Package) IsValid() bool {
	if p.astPackage == nil {
		return false
	}

	if p.buildPackage == nil {
		return false
	}

	if p.docPackage == nil {
		return false
	}

	// All checks passed.
	return true
}

// Name returns the package's name.
func (p Package) Name() string {
	if !p.IsValid() {
		return ""
	}

	return p.docPackage.Name
}

// Files returns a list of source files in the package. The file paths are relative to the package's
// directory, not absolute on the filesystem. Test files (*_test.go) are not included in the list.
// To get a list of test files in the package, see Package's TestFiles. Note: This returns all
// source files in the package's directory and does not limit the files based on what is actually
// used when building for the current system.
func (p Package) Files() []string {
	if !p.IsValid() {
		return nil
	}

	return p.files
}

// TestFiles returns a list of test files in the package's directory. This includes test files both
// within the package (e.g. mypkg_test.go in package mypkg) and outside of the package but within
// the package's directory (e.g. other_test.go in package mypkg_test). The file paths are relative
// to the package's directory, not absolute on the filesystem. Source files are not included in the
// list. To get a list of source files in the package, see Package's Files.
func (p Package) TestFiles() []string {
	if !p.IsValid() {
		return nil
	}

	return p.testFiles
}

// Imports returns a list of imports in the package. The list includes only imports from the source
// files, not the test files. To get a list of imports from the test files, see Package's TestImports.
func (p Package) Imports() []string {
	if !p.IsValid() {
		return nil
	}

	return p.imports
}

// TestImports returns a list of imports from test files both within the package (e.g. mypkg_test.go
// in package mypkg) and outside of the package but within the package's directory (e.g.
// other_test.go in package mypkg_test). The list includes only imports from the test files, not the
// source files. To get a list of imports from the source files, see Package's Imports.
func (p Package) TestImports() []string {
	if !p.IsValid() {
		return nil
	}

	return p.testImports
}

// Functions returns a list of exported functions in the package. The list includes exported
// functions from source files for the package only, not from test files (internal or external).
func (p Package) Functions() []Function {
	if !p.IsValid() {
		return nil
	}

	return p.functions
}

// Types returns a list of exported types in the package.
func (p Package) Types() []Type {
	if !p.IsValid() {
		return nil
	}

	return p.types
}
