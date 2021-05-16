// Package pkg provides a generic convenience wrapper around various golang package libraries.
package pkg

import (
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
)

type Package struct {
	buildPackage *build.Package
	astPackage   *ast.Package
	docPackage   *doc.Package
}
