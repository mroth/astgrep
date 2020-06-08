package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"sync"

	"github.com/mroth/astgrep"
)

// search takes a list of files to search concurrently using all Matchers
// provided.
//
// For now, the fileParser segment of this will log directly to os.Stderr if
// there are any file read errors, in the future if we want to integrate this
// into the public API of astgrep as a convenience tool we should adjust it to
// either return an err chan or make an optional io.Writer for error logging.
func search(files []string, matchers []astgrep.Matcher) (*token.FileSet, <-chan astgrep.Match) {
	fset, parsedFiles := parseFiles(files)
	resC := matchFiles(parsedFiles, matchers)
	return fset, resC
}

func parseFiles(files []string) (*token.FileSet, <-chan *ast.File) {
	fset := token.NewFileSet()
	parsedFiles := make(chan *ast.File)
	go func() {
		defer close(parsedFiles)
		for _, fp := range files {
			f, err := parser.ParseFile(fset, fp, nil, parser.ParseComments)
			if err != nil {
				// TODO: no direct logging in the future if exposing this in API
				log.Println(err)
			} else {
				parsedFiles <- f
			}
		}
	}()
	return fset, parsedFiles
}

func matchFiles(pf <-chan *ast.File, ms []astgrep.Matcher) <-chan astgrep.Match {
	resC := make(chan astgrep.Match)
	go func() {
		defer close(resC)
		for f := range pf {
			// for now, we run multiple matchers concurrently, but locked to a
			// single file, e.g. all have to finish before we move on to the
			// next file. This is not necessarily the most efficient way to
			// handle things, as faster matchers may possibly "waste" time
			// waiting for the others to finish, but reduces complexity vs
			// dealing with a fanout queue and our concurrency is already
			// benchmarking good enough for most applications.
			var wg sync.WaitGroup
			for _, m := range ms {
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
	return resC
}
