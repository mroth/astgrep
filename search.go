package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"sync"
)

func Search(files []string, matchers []Matcher) (*token.FileSet, <-chan Match) {
	fset := token.NewFileSet()

	parsedFiles := make(chan *ast.File)
	go func() {
		defer close(parsedFiles)
		for _, fp := range files {
			f, err := parser.ParseFile(fset, fp, nil, parser.ParseComments)
			if err != nil {
				log.Println(err) // TODO need to bubble
			} else {
				parsedFiles <- f
			}
		}
	}()

	resC := make(chan Match)
	go func() {
		defer close(resC)
		for f := range parsedFiles {
			// for now, we run multiple matchers concurrently, but locked to a
			// single file, e.g. all have to finish before we move on to the
			// next file. This is not necessarily the most efficient way to
			// handle things, since faster matchers may "waste" time waiting for
			// the others to finish, but is simpler than dealing with a fanout
			// queue and our concurrency is already good enough for most
			// applications.
			var wg sync.WaitGroup
			for _, m := range matchers {
				m := m
				wg.Add(1)
				go func() {
					defer wg.Done()
					m.Reset()
					ast.Walk(m, f)
					for _, m := range m.Matches() {
						resC <- m
					}
				}()
			}
			wg.Wait()
		}
	}()

	return fset, resC
}
