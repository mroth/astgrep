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
	pattern = flag.String("pattern", ".*", "pattern to search for")
)

func usage() {
	msg := "usage: %s [options] <file ...>\n"
	fmt.Fprintf(flag.CommandLine.Output(), msg, os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
	}

	re, err := regexp.Compile(*pattern)
	if err != nil {
		log.Fatal(err)
	}

	files := flag.Args()
	fset := token.NewFileSet()
	for _, fp := range files {
		f, err := parser.ParseFile(fset, fp, nil, 0)
		if err != nil {
			log.Println(err)
			continue
		}
		v := visitor{fset: fset, re: re}
		ast.Walk(v, f)
	}
}

type visitor struct {
	fset *token.FileSet
	re   *regexp.Regexp
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	// do not parse import specifications
	if _, ok := n.(*ast.ImportSpec); ok {
		return nil
	}

	if bl, ok := n.(*ast.BasicLit); ok && bl.Kind == token.STRING {
		position := v.fset.Position(bl.Pos())
		if v.re.MatchString(bl.Value) {
			fmt.Printf("%v\t%v\n", position, bl.Value)
		}
	}

	return v
}
