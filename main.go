package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
)

var (
	strPattern     = flag.String("string", "", "find string literals containing `pattern`")
	commentPattern = flag.String("comment", "", "find comments containing `pattern`")
	varPattern     = flag.String("var", "", "find variables or constants with name containing `pattern`")
	// numWorkers = flag.Int("workers", runtime.NumCPU(), "number of search threads")
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

	var matchers []Matcher
	if *strPattern != "" {
		re, err := regexp.Compile(*strPattern)
		if err != nil {
			log.Fatal(err)
		}
		matchers = append(matchers, &StrPatternVisitor{re: re})
	}

	if *commentPattern != "" {
		re, err := regexp.Compile(*commentPattern)
		if err != nil {
			log.Fatal(err)
		}
		matchers = append(matchers, &CommentPatternVisitor{re: re})
	}

	if *varPattern != "" {
		re, err := regexp.Compile(*varPattern)
		if err != nil {
			log.Fatal(err)
		}
		matchers = append(matchers, &VarPatternVisitor{re: re})
	}

	fset, resC := Search(files, matchers)
	for m := range resC {
		position := fset.Position(m.Pos())
		// fmt.Printf("%v\t%v\n", position, m.Text)
		fmt.Printf("%v\t%s%s%s%s%s\n", position,
			m.Text[:m.Base],
			"\033[31m", //red
			m.Text[m.Base:m.Base+m.Length],
			"\033[0m", //reset
			m.Text[m.Base+m.Length:],
		)
	}
}
