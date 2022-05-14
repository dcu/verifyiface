package analyzer

import (
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	verbose = true

	path, err := filepath.Abs("../samples")
	if err != nil {
		t.Fail()
	}

	_ = analysistest.Run(t, path, Analyzer, "example1", "example2", "example3", "example4", "example5", "example6")
}
