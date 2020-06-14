package main

import (
	"github.com/tdakkota/schemaverify"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(schemaverify.NewAnalyzer())
}
