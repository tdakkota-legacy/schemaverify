package schemaverify

import (
	"fmt"
	"go/ast"
	"path/filepath"

	"github.com/tdakkota/schemaverify/object"

	"github.com/lestrrat-go/jsschema"
	"golang.org/x/tools/go/analysis"
)

type SchemaVerify struct {
	SchemaDir     *string
	ReportSkipped *bool
}

func NewAnalyzer() *analysis.Analyzer {
	analyzer := &analysis.Analyzer{
		Name: "schemaverify",
		Doc:  "checks that generated Go structs matches JSON Schema",
	}

	schemaDir := analyzer.Flags.String("schema-dir", "vk-api-schema", "path to JSON Schema directory")
	reportSkipped := analyzer.Flags.Bool("report-skipped", false, "report definition which not found")

	s := SchemaVerify{SchemaDir: schemaDir, ReportSkipped: reportSkipped}
	analyzer.Run = s.Run

	return analyzer
}

func (s SchemaVerify) Run(pass *analysis.Pass) (interface{}, error) {
	if s.SchemaDir == nil {
		return nil, fmt.Errorf("schema-dir is empty")
	}

	objectsSchema, err := schema.ReadFile(filepath.Join(*s.SchemaDir, "objects.json"))
	if err != nil {
		return nil, err
	}

	return s.analyze(pass, objectsSchema)
}

func (s SchemaVerify) analyze(pass *analysis.Pass, sch *schema.Schema) (interface{}, error) {
	objects := object.NewMapping(sch)

	for _, file := range pass.Files {
		ast.Inspect(file, objects.InspectCallback)
	}

	return object.NewVerifier(objects, *s.ReportSkipped).Verify(pass)
}
