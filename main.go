package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"regexp"
)

var (
	strPattern = flag.String("string", ".*", "find string literals containing `pattern`")
	// commentPattern = flag.String("comment", "", "find comments containing `pattern`")
	// varPattern     = flag.String("var", "", "find variables containing `pattern`")
)

func usage() {
	msg := "usage: %s [options] <file ...>\n"
	fmt.Fprintf(flag.CommandLine.Output(), msg, os.Args[0])
	flag.PrintDefaults()
}

func main() {
	// parse CLI args
	flag.Usage = usage
	flag.Parse()
	files := flag.Args()
	if len(files) < 1 {
		flag.Usage()
	}

	// compile user-provided pattern
	re, err := regexp.Compile(*strPattern)
	if err != nil {
		log.Fatal(err)
	}

	fset := token.NewFileSet()
	var results []*ast.BasicLit
	for _, fp := range files {
		matches, err := parseFile(fset, fp, re)
		if err != nil {
			log.Println(err)
		} else {
			results = append(results, matches...)
		}
	}

	for _, m := range results {
		position := fset.Position(m.Pos())
		fmt.Printf("%v\t%v\n", position, m.Value)
	}
}

func parseFile(fset *token.FileSet, path string, re *regexp.Regexp) ([]*ast.BasicLit, error) {
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}
	v := strPatternVisitor{re: re}
	ast.Walk(&v, f)
	return v.matches, nil
}
