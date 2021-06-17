// Package pkg provides a generic convenience wrapper around various golang package libraries.
package pkg

import (
	"fmt"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
)

var (
	ErrInvalidPkg = fmt.Errorf("invalid package")
)

// Package is the main type for this package. It holds details about the package.
type Package struct {
	// Package object from go/build.
	buildPackage *build.Package

	// Package object from go/doc.
	docPackage *doc.Package
}

// New parses the package at importPath and creates a new Package object with its information.
func New(importPath string) (Package, error) {
	buildPackage, err := build.Import(importPath, "", 0)
	if err != nil {
		return Package{}, err
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, buildPackage.Dir, nil, parser.AllErrors)
	if err != nil {
		return Package{}, err
	}

	astPackage, ok := pkgs[buildPackage.Name]
	if !ok {
		return Package{}, fmt.Errorf("package not found in %s", importPath)
	}

	docPackage := doc.New(astPackage, importPath, 0)
	if docPackage == nil {
		return Package{}, ErrInvalidPkg
	}

	p := Package{
		buildPackage: buildPackage,
		docPackage:   docPackage,
	}

	return p, nil
}

// isValid checks whether or not the package object has valid data.
func (p Package) isValid() bool {
	if p == (Package{}) {
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
	if !p.isValid() {
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
	if !p.isValid() {
		return nil
	}

	// Return both the source files for this system's build and the source files ignored for this
	// system's build.
	return append(p.buildPackage.GoFiles, p.buildPackage.IgnoredGoFiles...)
}

// TestFiles returns a list of test files in the package's directory. This includes test files both
// within the package (e.g. mypkg_test.go in package mypkg) and outside of the package (e.g.
// other_test.go in package mypkg_test) but within the package's directory. The file paths are
// relative to the package's directory, not absolute on the filesystem. Source files are not
// included in the list. To get a list of source files in the package, see Package's Files.
func (p Package) TestFiles() []string {
	if !p.isValid() {
		return nil
	}

	return append(p.buildPackage.TestGoFiles, p.buildPackage.XTestGoFiles...)
}

// Imports returns a list of imports in the package. The list includes only imports from the source
// files, not the test files. To get a list of imports from the test files, see Package's TestImports.
func (p Package) Imports() []string {
	if !p.isValid() {
		return nil
	}

	return p.buildPackage.Imports
}

// TestImports returns a list of imports in the package. The list includes only imports from the
// test files, not the source files. To get a list of imports from the source files, see Package's
// Imports.
func (p Package) TestImports() []string {
	if !p.isValid() {
		return nil
	}

	return p.buildPackage.TestImports
}

// Types returns a list of exported types in the package.
func (p Package) Types() []Type {
	if !p.isValid() {
		return nil
	}

	// If there aren't any exported types in this package, then don't return anything.
	if len(p.docPackage.Types) == 0 {
		return nil
	}

	// Wrap every go/doc Type in our own Type.
	types := make([]Type, len(p.docPackage.Types))
	for i, v := range p.docPackage.Types {
		types[i] = Type{
			docType: v,
		}
	}

	return types
}

// Functions returns a list of exported functions in the package.
func (p Package) Functions() []Function {
	if !p.isValid() {
		return nil
	}

	// If there aren't any exported functions in this package, then don't return anything.
	if len(p.docPackage.Types) == 0 {
		return nil
	}

	// Wrap every go/doc Func in our own Function.
	funcs := make([]Function, len(p.docPackage.Funcs))
	for i, v := range p.docPackage.Funcs {
		funcs[i] = Function{
			docFunc: v,
		}
	}

	return funcs
}
