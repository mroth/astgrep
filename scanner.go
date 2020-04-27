package main

import (
	"go/ast"
	"go/token"
	"regexp"
)

// strPatternVisitor is an ast.Visitor that collects any ast.BasicLit node
// matching pattern re.
type strPatternVisitor struct {
	re      *regexp.Regexp
	matches []*ast.BasicLit // all will be of kind token.STRING
}

func (v *strPatternVisitor) Visit(n ast.Node) ast.Visitor {
	if _, ok := n.(*ast.ImportSpec); ok {
		return nil
	}

	if bl, ok := n.(*ast.BasicLit); ok && bl.Kind == token.STRING {
		if v.re.MatchString(bl.Value) {
			v.matches = append(v.matches, bl)
		}
	}

	return v
}
