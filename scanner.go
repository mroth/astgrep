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

// StrPatternVisitor is a Matcher that finds all matches of a given pattern
// matching against the Value of ast.BasicList nodes of type token.STRING.
type StrPatternVisitor struct {
	re      *regexp.Regexp
	matches []Match
}

// Visit implements ast.Visitor, it is thus used during ast.Walk and typically
// will not be called directly.
func (v *StrPatternVisitor) Visit(n ast.Node) ast.Visitor {
	if _, ok := n.(*ast.ImportSpec); ok {
		// skip import specs (technically contain string literals, but no one
		// cares.)
		return nil
	}

	if bl, ok := n.(*ast.BasicLit); ok && bl.Kind == token.STRING {
		s := bl.Value
		if locs := v.re.FindAllStringIndex(s, -1); locs != nil {
			for _, loc := range locs {
				v.matches = append(v.matches, Match{
					Node:   bl,
					Text:   s[loc[0]:loc[1]],
					Base:   loc[0],
					Length: loc[1] - loc[0],
				})
			}
		}
	}

	return v
}

// Matches returns all matches collected by the Matcher
func (v *StrPatternVisitor) Matches() []Match {
	return v.matches
}

// CommentPatternVisitor is a Matcher that finds all matches of a given pattern
// matching against the Text of ast.Comment nodes.
type CommentPatternVisitor struct {
	re      *regexp.Regexp
	matches []Match
}

// Visit implements ast.Visitor, it is thus used during ast.Walk and typically
// will not be called directly.
func (v *CommentPatternVisitor) Visit(n ast.Node) ast.Visitor {
	if c, ok := n.(*ast.Comment); ok {
		s := c.Text
		if locs := v.re.FindAllStringIndex(s, -1); locs != nil {
			for _, loc := range locs {
				v.matches = append(v.matches, Match{
					Node:   c,
					Text:   s[loc[0]:loc[1]],
					Base:   loc[0],
					Length: loc[1] - loc[0],
				})
			}
		}
	}
	return v
}

// Matches returns all matches collected by the Matcher
func (v *CommentPatternVisitor) Matches() []Match {
	return v.matches
}
