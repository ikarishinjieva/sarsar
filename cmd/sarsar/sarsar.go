package main

import (
	"github.com/ikarishinjieva/sarsar/sarsar"
	"flag"
	"fmt"
	"os"
)

var fInputFile string
var fHelp bool

func init() {
	flag.StringVar(&fInputFile, "f", "", "input file")
	flag.BoolVar(&fHelp, "h", false, "print help message")
}

func main() {
	flag.Parse()

	if fHelp {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

	if "" == fInputFile {
		fmt.Fprint(os.Stderr, "Error: input file required\n")
		os.Exit(1)
	}

	if err := sarsar.SarSar(fInputFile); nil != err {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}

