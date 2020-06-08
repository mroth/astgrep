package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/gookit/color"
	"github.com/mroth/astgrep"
)

var (
	strPattern     = flag.String("string", "", "find string literals containing `pattern`")
	commentPattern = flag.String("comment", "", "find comments containing `pattern`")
	varPattern     = flag.String("var", "", "find variables or constants with name containing `pattern`")
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

	var matchers []astgrep.Matcher
	if *strPattern != "" {
		re, err := regexp.Compile(*strPattern)
		if err != nil {
			log.Fatal(err)
		}
		matchers = append(matchers, astgrep.NewStrPatternMatcher(re))
	}

	if *commentPattern != "" {
		re, err := regexp.Compile(*commentPattern)
		if err != nil {
			log.Fatal(err)
		}
		matchers = append(matchers, astgrep.NewCommentPatternMatcher(re))
	}

	if *varPattern != "" {
		re, err := regexp.Compile(*varPattern)
		if err != nil {
			log.Fatal(err)
		}
		matchers = append(matchers, astgrep.NewVarPatternMatcher(re))
	}

	fset, resC := search(files, matchers)
	for m := range resC {
		position := fset.Position(m.Pos())
		fmt.Printf("%v\t%s%s%s\n",
			position,
			m.Text[:m.Base],
			color.FgRed.Render(m.Text[m.Base:m.Base+m.Length]), // matched portion
			m.Text[m.Base+m.Length:],
		)
	}
}
