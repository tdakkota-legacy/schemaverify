package analyzer

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"log"
	"path/filepath"
	"testing"
)

func TestAnalyzer(t *testing.T) {
	testdata, err := filepath.Abs("../testdata/objects")
	if err != nil {
		log.Fatal(err)
	}

	analyzer := NewAnalyzer()
	err = analyzer.Flags.Set("schema-dir", "../testdata/schema")
	if err != nil {
		log.Fatal(err)
	}

	analysistest.Run(t, testdata, analyzer)
}
