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
	strPattern = flag.String("string", ".*", "find string literals containing `pattern`")
	// commentPattern = flag.String("comment", "", "find comments containing `pattern`")
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

	// compile user-provided pattern
	re, err := regexp.Compile(*strPattern)
	if err != nil {
		log.Fatal(err)
	}

	fset := token.NewFileSet()
	// parseErrs := make(chan error) dont want to block on error reading
	parsedFiles := make(chan *ast.File)
	go func() {
		defer close(parsedFiles)
		for _, fp := range files {
			f, err := parser.ParseFile(fset, fp, nil, 0)
			if err != nil {
				log.Println(err)
			} else {
				parsedFiles <- f
			}
		}
	}()

	// TODO while excess workers...

	// str scanner workers
	var wg sync.WaitGroup
	resC := make(chan []*ast.BasicLit)
	for i := 0; i < *numWorkers; i++ {
		wg.Add(1)
		go func() {
			v := strPatternVisitor{re: re}
			for f := range parsedFiles {
				v.matches = nil // realloc slice to reset
				ast.Walk(&v, f)
				resC <- v.matches
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(resC)
	}()

	for r := range resC {
		for _, m := range r {
			position := fset.Position(m.Pos())
			fmt.Printf("%v\t%v\n", position, m.Value)
		}
	}
}
