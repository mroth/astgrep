package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"testing"
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

func TestStrPatternVisitor(t *testing.T) {
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
		// TODO: .*
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			f, _ := parseTestFile(t, tt.path)
			v := StrPatternVisitor{re: tt.re}
			ast.Walk(&v, f)
			if actualMatches := len(v.matches); actualMatches != tt.numMatches {
				t.Errorf(
					"%v %v: want %d matches, got %d",
					tt.path, tt.re, tt.numMatches, actualMatches,
				)
			}
		})
	}
}

func TestCommentPatternVisitor(t *testing.T) {
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
		// TODO: .*
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			f, _ := parseTestFile(t, tt.path)
			v := CommentPatternVisitor{re: tt.re}
			ast.Walk(&v, f)
			if actualMatches := len(v.matches); actualMatches != tt.numMatches {
				t.Errorf(
					"%v %v: want %d matches, got %d",
					tt.path, tt.re, tt.numMatches, actualMatches,
				)
			}
		})
	}
}

func TestVarPatternVisitor(t *testing.T) {
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
			v := VarPatternVisitor{re: tt.re}
			ast.Walk(&v, f)
			if actualMatches := len(v.matches); actualMatches != tt.numMatches {
				t.Errorf(
					"%v %v: want %d matches, got %d",
					tt.path, tt.re, tt.numMatches, actualMatches,
				)
			}
		})
	}
}
