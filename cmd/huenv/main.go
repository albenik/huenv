package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/albenik/huenv/generator"
)

func main() {
	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	f.Usage = func() {
		w := os.Stderr
		fmt.Fprintln(w, "usage: huenv [flags] <package> <type>")
		fmt.Fprintln(w, "flags:")
		f.PrintDefaults()
	}

	out := f.String("out", "", "Destination file. Defaults to stdout.")
	codeOnly := f.Bool("code_only", false, "Only generate the reflection program, write it to stdout and exit.")
	buildFlagsStr := f.String("build_flags", "", "Additional flags for go build.")

	if err := f.Parse(os.Args[1:]); err != nil {
		f.Usage()
		printErrorAndExit(err, 2)
	}

	if f.NArg() != 2 {
		f.Usage()
		os.Exit(2)
	}

	srcPackage := f.Arg(0)
	srcType := f.Arg(1)

	if *codeOnly {
		src, err := generator.GenerateProgram(srcPackage, srcType)
		if err != nil {
			printErrorAndExit(err, 1)
		}
		fmt.Println(src)
		return
	}

	var buildFlags []string
	if *buildFlagsStr != "" {
		buildFlags = strings.Split(*buildFlagsStr, " ")
	}
	if err := generator.BuildAndRun(*out, srcPackage, srcType, buildFlags); err != nil {
		printErrorAndExit(err, 1)
	}
}

func printErrorAndExit(err error, code int) {
	fmt.Fprintln(os.Stderr, "huenv:", err)

	var derr generator.DetailedError
	if errors.As(err, &derr) {
		fmt.Fprint(os.Stderr, derr.Details())
	} else {
		fmt.Fprintln(os.Stderr, "no error details")
	}

	os.Exit(code)
}
