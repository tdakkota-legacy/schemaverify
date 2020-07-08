package schemaverify

import (
	"log"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata, err := filepath.Abs("./testdata/objects")
	if err != nil {
		log.Fatal(err)
	}

	analyzer := NewAnalyzer()
	err = analyzer.Flags.Set("schema-dir", "./testdata/schema")
	if err != nil {
		log.Fatal(err)
	}

	analysistest.Run(t, testdata, analyzer)
}
