package astgrep_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"testing"

	"github.com/mroth/astgrep"
)

func parseTestFile(t *testing.T, path string) (*ast.File, *token.FileSet) {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	return f, fset
}

func TestStrPatternMatcher(t *testing.T) {
	tests := []struct {
		path       string
		re         *regexp.Regexp
		numMatches int
	}{
		{
			path:       "testdata/sample.go",
			re:         regexp.MustCompile("apple"),
			numMatches: 1,
		},
		{
			path:       "testdata/sample.go",
			re:         regexp.MustCompile("alpha"),
			numMatches: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			f, _ := parseTestFile(t, tt.path)
			v := astgrep.NewStrPatternMatcher(tt.re)
			ast.Walk(v, f)
			if actualMatches := len(v.Matches()); actualMatches != tt.numMatches {
				t.Errorf(
					"%v %v: want %d matches, got %d",
					tt.path, tt.re, tt.numMatches, actualMatches,
				)
			}
		})
	}
}

func TestCommentPatternMatcher(t *testing.T) {
	tests := []struct {
		path       string
		re         *regexp.Regexp
		numMatches int
	}{
		{
			path:       "testdata/sample.go",
			re:         regexp.MustCompile("banana"),
			numMatches: 1,
		},
		{
			path:       "testdata/sample.go",
			re:         regexp.MustCompile("cucmber"),
			numMatches: 0,
		},
		{
			path:       "testdata/sample.go",
			re:         regexp.MustCompile("alpha"),
			numMatches: 1,
		},
		{
			path:       "testdata/sample.go",
			re:         regexp.MustCompile("(?i)alpha"),
			numMatches: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			f, _ := parseTestFile(t, tt.path)
			v := astgrep.NewCommentPatternMatcher(tt.re)
			ast.Walk(v, f)
			if actualMatches := len(v.Matches()); actualMatches != tt.numMatches {
				t.Errorf(
					"%v %v: want %d matches, got %d",
					tt.path, tt.re, tt.numMatches, actualMatches,
				)
			}
		})
	}
}

func TestVarPatternMatcher(t *testing.T) {
	tests := []struct {
		path       string
		re         *regexp.Regexp
		numMatches int
	}{
		{
			path:       "testdata/sample.go",
			re:         regexp.MustCompile("alpha"),
			numMatches: 1, // dont match param names
		},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			f, _ := parseTestFile(t, tt.path)
			v := astgrep.NewVarPatternMatcher(tt.re)
			ast.Walk(v, f)
			if actualMatches := len(v.Matches()); actualMatches != tt.numMatches {
				t.Errorf(
					"%v %v: want %d matches, got %d",
					tt.path, tt.re, tt.numMatches, actualMatches,
				)
			}
		})
	}
}
