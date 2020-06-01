package main

import (
	"go/ast"
	"go/token"
	"regexp"
)

// Matcher is an ast.Visitor which collects all Match based on a specific
// matching criteria during its Walk.
type Matcher interface {
	Visit(n ast.Node) ast.Visitor // ast.Visitor implementation
	Matches() []Match             // return existing matches
	Reset()                       // clear existing matches
}

// Match represents a positive result found by a Matcher.
type Match struct {
	Node   ast.Node // node matched against
	Text   string   // contents used for match search
	Base   int      // base offset of match within text
	Length int      // length of match text
}

// Pos returns the position of first character belonging to the Match.
func (m Match) Pos() token.Pos {
	return token.Pos(int(m.Node.Pos()) + m.Base)
}

// End returns the position of first character immediately after the Match.
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

// StrPatternMatcher is a Matcher that finds all matches of a given pattern
// matching against the Value of ast.BasicList nodes of type token.STRING.
type StrPatternMatcher struct {
	re      *regexp.Regexp
	matches []Match
}

// Visit implements ast.Visitor, it is thus used during ast.Walk and typically
// will not be called directly.
func (v *StrPatternMatcher) Visit(n ast.Node) ast.Visitor {
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
					Text:   s,
					Base:   loc[0],
					Length: loc[1] - loc[0],
				})
			}
		}
	}

	return v
}

// Matches returns all matches collected by the Matcher.
func (v *StrPatternMatcher) Matches() []Match {
	return v.matches
}

// Reset the collected Matches to be empty.
func (v *StrPatternMatcher) Reset() {
	v.matches = nil
}

// CommentPatternMatcher is a Matcher that finds all matches of a given pattern
// matching against the Text of ast.Comment nodes.
type CommentPatternMatcher struct {
	re      *regexp.Regexp
	matches []Match
}

// Visit implements ast.Visitor, it is thus used during ast.Walk and typically
// will not be called directly.
func (v *CommentPatternMatcher) Visit(n ast.Node) ast.Visitor {
	if c, ok := n.(*ast.Comment); ok {
		s := c.Text
		if locs := v.re.FindAllStringIndex(s, -1); locs != nil {
			for _, loc := range locs {
				v.matches = append(v.matches, Match{
					Node:   c,
					Text:   s,
					Base:   loc[0],
					Length: loc[1] - loc[0],
				})
			}
		}
	}
	return v
}

// Matches returns all matches collected by the Matcher.
func (v *CommentPatternMatcher) Matches() []Match {
	return v.matches
}

// Reset the collected Matches to be empty.
func (v *CommentPatternMatcher) Reset() {
	v.matches = nil
}

// VarPatternMatcher is a Matcher that finds all matches of a given pattern
// matching against a variable of constant declaration (the Names of
// ast.ValueSpec node).
type VarPatternMatcher struct {
	re      *regexp.Regexp
	matches []Match
}

// Visit implements ast.Visitor, it is thus used during ast.Walk and typically
// will not be called directly.
func (v *VarPatternMatcher) Visit(n ast.Node) ast.Visitor {
	if c, ok := n.(*ast.ValueSpec); ok {
		s := c.Names[0].Name // TODO: what situation is len(Names) > 1?
		if locs := v.re.FindAllStringIndex(s, -1); locs != nil {
			for _, loc := range locs {
				v.matches = append(v.matches, Match{
					Node:   c,
					Text:   s,
					Base:   loc[0],
					Length: loc[1] - loc[0],
				})
			}
		}
	}
	return v
}

// Matches returns all matches collected by the Matcher.
func (v *VarPatternMatcher) Matches() []Match {
	return v.matches
}

// Reset the collected Matches to be empty.
func (v *VarPatternMatcher) Reset() {
	v.matches = nil
}
