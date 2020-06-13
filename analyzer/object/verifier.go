package object

import (
	"fmt"
	"go/ast"
	"schemaverify/analyzer/parseutil"
	"schemaverify/analyzer/pragma"

	schema "github.com/lestrrat-go/jsschema"
	"golang.org/x/tools/go/analysis"
)

type Verifier struct {
	Objects       SchemaObjects
	Schema        *schema.Schema
	ReportSkipped bool
}

func NewVerifier(objects SchemaObjects, schema *schema.Schema, reportSkipped bool) Verifier {
	return Verifier{Objects: objects, Schema: schema, ReportSkipped: reportSkipped}
}

func (o Verifier) Verify(pass *analysis.Pass) (interface{}, error) {
	err := o.Objects.ForEach(func(structName string, alias []string, object SchemaObject) (bool, error) {
		def, ok := Definitions(o.Schema.Definitions).Find(structName, alias)
		if !ok { // definition not found: skip
			if o.ReportSkipped {
				pass.Reportf(
					object.Name.Pos(),
					"%s (%s) struct not found in schema",
					structName,
					parseutil.PascalToSnake(structName),
				)
			}
			return true, nil
		}

		err := o.verifyObject(pass, object, def)
		if err != nil {
			return false, err
		}

		return true, nil
	})

	return nil, err
}

const debug = false

func (o Verifier) verifyObject(pass *analysis.Pass, obj SchemaObject, def *schema.Schema) error {
	typeName := obj.Name.Name

	if debug {
		fmt.Println(typeName, ":", parseutil.PascalToSnake(typeName))
	}

	switch v := obj.Type.(type) {
	case *ast.StructType:
		return o.verifyStruct(pass, typeName, v, def)
	}

	return nil
}

func (o Verifier) verifyStruct(
	pass *analysis.Pass,
	typeName string,
	obj *ast.StructType,
	def *schema.Schema,
) error {
	fields := o.Objects.MapFields(obj.Fields.List)

	for _, name := range def.Required {
		if _, ok := fields[name]; !ok {
			pass.Reportf(
				obj.Struct,
				"%s field is required in object %s (%s)",
				name, typeName, parseutil.PascalToSnake(typeName),
			)
		}
	}

	for name, field := range fields {
		pragmas := pragma.ParsePragmas(field.Doc)
		if pragmas.Skip() {
			continue
		}

		prop := def.Properties[name]
		if prop == nil {
			if o.ReportSkipped {
				pass.Reportf(
					obj.Struct,
					"%s (%s) field not found in schema of `%s` object",
					field.Names[0], name, parseutil.PascalToSnake(typeName),
				)
			}
			continue
		}

		schemaTypeName, err := o.verifyField(pass.TypesInfo, field.Type, prop)
		if err != nil {
			pass.Reportf(field.Pos(), "Field %s does not match schema: %s", field.Names[0], err.Error())
		}

		if debug {
			fmt.Printf("\t%v(%s): %v", field.Names, name, field.Type)
			fmt.Print(" -> ", schemaTypeName)
			fmt.Println()
		}
	}

	return nil
}
