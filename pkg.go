// Package pkg provides a generic convenience wrapper around various libraries in the standard
// library.
package pkg

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strings"
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

	// List of directories within this package's directory.
	subdirectories []string

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
	// Generate the go/build Package for the import path so we can have more visibility into this
	// package's structure.
	buildPkg, err := build.Import(importPath, "", 0)
	if err != nil {
		return Package{}, fmt.Errorf("build error: invalid package in %s: %w", importPath, err)
	}

	fset := token.NewFileSet()
	astPkgs, err := parser.ParseDir(fset, buildPkg.Dir, nil, parser.ParseComments)
	if err != nil {
		return Package{}, fmt.Errorf("invalid package in %s: %w", importPath, err)
	}

	// Get the go/ast Package for the package named by the import path.
	astPkg, ok := astPkgs[buildPkg.Name]
	if !ok {
		return Package{}, fmt.Errorf("package not found in %s", importPath)
	}
	if astPkg == nil {
		return Package{}, fmt.Errorf("missing package %s in %s", importPath, buildPkg.Dir)
	}

	// Generate the go/doc Package for the package named by the import path. We first have to
	// flatten out the map of ast files and then use that list to parse the individual files. We
	// want to gather all files for the package and not filter out any based on build systems.
	astFiles := make([]*ast.File, 0, len(astPkg.Files))
	for _, v := range astPkg.Files {
		astFiles = append(astFiles, v)
	}
	docPkg, err := doc.NewFromFiles(fset, astFiles, importPath)
	if err != nil {
		return Package{}, fmt.Errorf("invalid package in %s: %w", importPath, err)
	}
	if docPkg == nil {
		return Package{}, ErrInvalidPkg
	}

	// Put everything together into our Package type.
	return newPackage(buildPkg, docPkg, fset)
}

// newPackage puts together the internal structure for a Package object.
func newPackage(buildPkg *build.Package, docPkg *doc.Package, fset *token.FileSet) (Package, error) {
	// Begin with structuring up our object with what we have so far.
	pkg := Package{
		name:       docPkg.Name,
		importPath: docPkg.ImportPath,
		comments:   docPkg.Doc,
	}

	// Put together the list of source files, both for this system's build and those ignored for this
	// system's build. We have to manually exclude test files because they might be in the same package.
	for _, ss := range [][]string{buildPkg.GoFiles, buildPkg.IgnoredGoFiles} {
		for _, s := range ss {
			if !strings.HasSuffix(s, "_test.go") {
				pkg.files = append(pkg.files, s)
			}
		}
	}
	sort.Strings(pkg.files)

	// Put together the list of test files, both for this package and any other external test package
	// in this package's directory. We have to manually add test files that are in the main package.
	for _, ss := range [][]string{buildPkg.TestGoFiles, buildPkg.XTestGoFiles} {
		pkg.testFiles = append(pkg.testFiles, ss...)
	}
	for _, ss := range [][]string{buildPkg.GoFiles, buildPkg.IgnoredGoFiles} {
		for _, s := range ss {
			if strings.HasSuffix(s, "_test.go") {
				pkg.testFiles = append(pkg.testFiles, s)
			}
		}
	}
	sort.Strings(pkg.testFiles)

	// Find all the subdirectories within this package's directory.
	subs, _ := os.ReadDir(buildPkg.Dir)
	for _, sub := range subs {
		if sub.IsDir() {
			pkg.subdirectories = append(pkg.subdirectories, sub.Name())
		}
	}
	sort.Strings(pkg.subdirectories)

	// Copy the list of imports from the source files.
	pkg.imports = append([]string{}, buildPkg.Imports...)

	// Put together the list of imports from both the internally and externally packaged test files.
	// Make sure the list is sorted and free of duplicates.
	m := make(map[string]bool)
	for _, ss := range [][]string{buildPkg.TestImports, buildPkg.XTestImports} {
		for _, s := range ss {
			m[s] = true
		}
	}
	pkg.testImports = make([]string, 0, len(m))
	for v := range m {
		pkg.testImports = append(pkg.testImports, v)
	}
	sort.Strings(pkg.testImports)

	// Extract the blocks of exported constants for this package, both for standard types (go/doc's
	// Consts) and for custom types (go/doc's Type's Consts).
	for _, cb := range docPkg.Consts {
		pkg.constantBlocks = append(pkg.constantBlocks, newConstantBlock(cb, nil, fset))
	}
	for _, t := range docPkg.Types {
		for _, cb := range t.Consts {
			pkg.constantBlocks = append(pkg.constantBlocks, newConstantBlock(cb, t, fset))
		}
	}

	// Extract the blocks of exported variables for this package, both for standard types (go/doc's
	// Vars) and for custom types (go/doc's Type's Vars).
	for _, vb := range docPkg.Vars {
		pkg.variableBlocks = append(pkg.variableBlocks, newVariableBlock(vb, nil, fset))
	}
	for _, t := range docPkg.Types {
		for _, vb := range t.Vars {
			pkg.variableBlocks = append(pkg.variableBlocks, newVariableBlock(vb, t, fset))
		}
	}

	// Extract the exported functions for this package.
	pkg.functions = make([]Function, len(docPkg.Funcs))
	for i, f := range docPkg.Funcs {
		pkg.functions[i] = newFunction(f, fset)
	}

	// Extract the exported types for this package.
	pkg.types = make([]Type, len(docPkg.Types))
	for i, t := range docPkg.Types {
		pkg.types[i] = newType(t, fset)
	}

	return pkg, nil
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

// Subdirectories returns a list of all subdirectories within this package's directory. The file
// paths are relative to the package's directory, not absolute on the filesystem.
func (p Package) Subdirectories() []string {
	return append([]string{}, p.subdirectories...)
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
