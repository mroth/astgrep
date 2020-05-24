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
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	return f, fset
}

func Test_strPatternVisitor(t *testing.T) {
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
			v := strPatternVisitor{re: tt.re}
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
