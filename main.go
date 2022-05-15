package main

import (
	"github.com/dcu/verifyiface/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	analyzer.Analyzer.Flags.BoolVar(&analyzer.Verbose, "verbose", false, "enable verbose mode")
	analyzer.Analyzer.Flags.BoolVar(&analyzer.StrictCheck, "strict", false, "enable strict mode")

	singlechecker.Main(analyzer.Analyzer)
}
