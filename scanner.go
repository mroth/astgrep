package main

import (
	"go/ast"
	"go/token"
	"regexp"
)

// Matcher is an ast.Visitor which collects all Match based on a specific
// matching criteria during its Walk.
type Matcher interface { // TODO better name?
	Visit(n ast.Node) ast.Visitor
	Matches() []Match
}

type Match struct {
	Node   ast.Node // node matched against
	Text   string   // contents used for match search
	Base   int      // base offset of match within text
	Length int      // length of match text
}

func (m Match) Pos() token.Pos {
	return token.Pos(int(m.Node.Pos()) + m.Base)
}

func (m Match) End() token.Pos {
	return token.Pos(int(m.Node.Pos()) + m.Base + m.Length)
}

/* Alternative where Text is accessed dynamically, but difficult to do in a
fully type safe way, and prevents external extensibility without modifying this
method. */

// Text returns the textual content the Match was peformed against. func (m
// *Match) Text() string {
//  node := *m.node
//  switch n := node.(type) {
//  case *ast.BasicLit:
//      return n.Value
//  case *ast.Comment:
//      return n.Text
//  default:
//      // should be unreachable, but cant do with type safety in Go since no
//      // exhaustiveness checking on switch statements nor real enum support.
//      panic("FATAL - 'unreachable' condition: unsupported ast.Node type")
//  }
// }

// strPatternVisitor is an ast.Visitor that collects any ast.BasicLit node
// matching pattern re.
type strPatternVisitor struct {
	re *regexp.Regexp
	// matches []*ast.BasicLit // all will be of kind token.STRING
	matches []Match
}

func (v *strPatternVisitor) Visit(n ast.Node) ast.Visitor {
	if _, ok := n.(*ast.ImportSpec); ok {
		// skip import specs (technically contain string literals, but no one
		// cares.)
		return nil
	}

	if bl, ok := n.(*ast.BasicLit); ok && bl.Kind == token.STRING {
		s := bl.Value
		if loc := v.re.FindStringIndex(s); loc != nil {
			v.matches = append(v.matches, Match{
				Node:   bl,
				Text:   s[loc[0]:loc[1]],
				Base:   loc[0],
				Length: loc[1] - loc[0],
			})
		}
	}

	return v
}

type commentPatternVisitor struct {
	re      *regexp.Regexp
	matches []Match
}

func (v *commentPatternVisitor) Visit(n ast.Node) ast.Visitor {
	if c, ok := n.(*ast.Comment); ok {
		s := c.Text
		if loc := v.re.FindStringIndex(s); loc != nil {
			v.matches = append(v.matches, Match{
				Node:   c,
				Text:   s[loc[0]:loc[1]],
				Base:   loc[0],
				Length: loc[1] - loc[0],
			})
		}
	}
	return v
}
