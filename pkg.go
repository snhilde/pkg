// Package pkg provides a generic convenience wrapper around various golang package libraries.
package pkg

import (
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
)

var (
	InvalidPkg = fmt.Errorf("invalid package")
)

type Package struct {
	buildPackage *build.Package
	astPackage   *ast.Package
	docPackage   *doc.Package
}

func New(importPath string) (*Package, error) {
	buildPackage, err := build.Import(importPath, "", 0)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, p.Dir, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	astPackage, ok := pkgs[buildPackage.Name]
	if !ok {
		return nil, fmt.Errorf("package not found in %s", importPath)
	}

	docPackage := doc.New(astPackage, importPath, 0)
	if doc == nil {
		return nil, InvalidPkg
	}

	p := new(Package)
	p.buildPackage = buildPackage
	p.astPackage   = astPackage
	p.docPackage   = docPackage

	return p, nil
}
