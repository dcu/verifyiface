package analyzer

import (
	"log"
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

	results := analysistest.Run(t, path, Analyzer, "example1")

	for _, r := range results {
		log.Printf("%#v", r)
	}
}
