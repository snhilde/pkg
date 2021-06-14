// Package pkg provides a generic convenience wrapper around various golang package libraries.
package pkg

import (
	"fmt"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"path/filepath"
)

var (
	InvalidPkg = fmt.Errorf("invalid package")
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
		return Package{}, InvalidPkg
	}

	p := Package {
		buildPackage: buildPackage,
		docPackage:   docPackage,
	}

	return p, nil
}

// valid checks whether or not the package object has valid data.
func (p Package) valid() bool {
	if p == (Package{}) {
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
	if !p.valid() {
		return ""
	}

	return p.docPackage.Name
}

// Files returns a list of files in the package. The file paths are relative to the package's
// directory, not absolute on the filesystem.
func (p Package) Files() []string {
	if !p.valid() {
		return nil
	}

	// go/doc's Package holds absolute paths for each file in the package. We want to convert those
	// to relative paths.
	basePath := p.buildPackage.Dir
	absPaths := p.docPackage.Filenames
	relPaths := make([]string, len(absPaths))
	for i, absPath := range absPaths {
		relPath, err := filepath.Rel(basePath, absPath)
		if err != nil {
			return nil
		}
		relPaths[i] = relPath
	}

	return relPaths
}

// Imports returns a list of imports in the package.
func (p Package) Imports() []string {
	if !p.valid() {
		return nil
	}

	return p.docPackage.Imports
}

// Types returns a list of exported types in the package.
func (p Package) Types() []Type {
	if !p.valid() {
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
	if !p.valid() {
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
