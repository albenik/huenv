package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/albenik/huenv/generator"
	"github.com/albenik/huenv/internal/version"
)

func main() {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.Usage = func() {
		w := os.Stderr
		fmt.Fprintln(w, "usage: huenv [flags] <package> <type>")
		fmt.Fprintln(w, "flags:")
		fs.PrintDefaults()
	}

	ver := fs.Bool("version", false, "Display version")
	out := fs.String("out", "", "Destination file. Defaults to stdout.")
	codeOnly := fs.Bool("code_only", false, "Only generate the reflection program, write it to stdout and exit.")
	buildFlagsStr := fs.String("build_flags", "", "Additional flags for go build.")

	if err := fs.Parse(os.Args[1:]); err != nil {
		fs.Usage()
		printErrorAndExit(err, 2) //nolint:gomnd
	}

	if *ver {
		fmt.Fprintln(os.Stdout, version.String())
		return
	}

	if fs.NArg() != 2 { //nolint:gomnd
		fs.Usage()
		os.Exit(2) //nolint:gomnd
	}

	srcPackage := fs.Arg(0)
	srcType := fs.Arg(1)

	if *codeOnly {
		src, err := generator.GenerateProgram(srcPackage, srcType)
		if err != nil {
			printErrorAndExit(err, 1)
		}
		fmt.Fprintln(os.Stdout, src)
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
