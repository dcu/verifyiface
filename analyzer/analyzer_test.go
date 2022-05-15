package analyzer

import (
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	Verbose = true

	path, err := filepath.Abs("../samples")
	if err != nil {
		t.Fail()
	}

	_ = analysistest.Run(t, path, Analyzer, "example1", "example2", "example3", "example4", "example5", "example6", "example7", "example8", "example9", "example10")
}
