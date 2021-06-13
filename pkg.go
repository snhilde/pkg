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
	InvalidPkg = fmt.Errorf("invalid package")
)

// Package is the main type for this package. It holds details about the package.
type Package struct {
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
		docPackage: docPackage,
	}

	return p, nil
}

// valid checks whether or not the package object has valid data.
func (p Package) valid() bool {
	if p == nil {
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
