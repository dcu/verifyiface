package main

import (
	"github.com/dcu/verifyiface/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	analyzer.Analyzer.Flags.BoolVar(&analyzer.Verbose, "verbose", false, "enable verbose mode")

	singlechecker.Main(analyzer.Analyzer)
}
