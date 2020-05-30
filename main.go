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
	"runtime"
	"sync"
)

var (
	strPattern     = flag.String("string", "", "find string literals containing `pattern`")
	commentPattern = flag.String("comment", "", "find comments containing `pattern`")
	// varPattern     = flag.String("var", "", "find variables containing `pattern`")
	numWorkers = flag.Int("workers", runtime.NumCPU(), "number of search threads")
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

	fset := token.NewFileSet()
	// parseErrs := make(chan error) dont want to block on error reading
	parsedFiles := make(chan *ast.File)
	go func() {
		defer close(parsedFiles)
		for _, fp := range files {
			f, err := parser.ParseFile(fset, fp, nil, parser.ParseComments)
			if err != nil {
				log.Println(err)
			} else {
				parsedFiles <- f
			}
		}
	}()

	var wg sync.WaitGroup
	resC := make(chan Match)

	// TODO while excess workers...
	// str scanner worker(+s later)
	if *strPattern != "" {
		re, err := regexp.Compile(*strPattern)
		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		go func() {
			v := StrPatternVisitor{re: re}
			for f := range parsedFiles {
				v.matches = nil // realloc slice to reset
				ast.Walk(&v, f)
				for _, m := range v.matches {
					resC <- m
				}
			}
			wg.Done()
		}()
	}

	// comment scanner worker(+s later)
	if *commentPattern != "" {
		re, err := regexp.Compile(*commentPattern)
		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		go func() {
			v := CommentPatternVisitor{re: re}
			for f := range parsedFiles {
				v.matches = nil // realloc slice to reset
				ast.Walk(&v, f)
				for _, m := range v.matches {
					resC <- m
				}
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(resC)
	}()

	for r := range resC {
		position := fset.Position(r.Pos())
		fmt.Printf("%v\t%v\n", position, r.Text)
	}
}
