package object

import (
	"fmt"
	"go/ast"
	"schemaverify/analyzer/pragma"

	schema "github.com/lestrrat-go/jsschema"
	"golang.org/x/tools/go/analysis"
)

type Verifier struct {
	Objects       SchemaObjects
	Schema        *schema.Schema
	ReportSkipped bool
}

func NewVerifier(objects SchemaObjects, reportSkipped bool) Verifier {
	return Verifier{Objects: objects, ReportSkipped: reportSkipped}
}

func (o Verifier) Verify(pass *analysis.Pass) (interface{}, error) {
	err := o.Objects.ForEach(func(structName string, object SchemaObject) (bool, error) {
		def, ok := o.Objects.FindDefinition(structName)
		if !ok { // definition not found: skip
			if o.ReportSkipped {
				pass.Reportf(
					object.Name.Pos(),
					"%s (%s) struct not found in schema",
					structName,
					def.SchemaName,
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

const debug = true

func (o Verifier) verifyObject(pass *analysis.Pass, obj SchemaObject, def DefinitionResult) error {
	typeName := obj.Name.Name

	if debug {
		fmt.Println(typeName, ":", def.SchemaName)
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
	def DefinitionResult,
) error {
	fields := o.Objects.MapFields(obj.Fields.List)
	props, err := MapProperties(def.Schema)
	if err != nil {
		return err
	}

	for _, name := range props.Required {
		if _, ok := fields[name]; !ok {
			pass.Reportf(
				obj.Struct,
				"%s field is required in object %s (%s)",
				name, typeName, def.SchemaName,
			)
		}
	}

	for name, field := range fields {
		pragmas := pragma.ParsePragmas(field.Doc)
		if pragmas.Skip() {
			continue
		}

		prop := props.Props[name]
		if prop == nil {
			if o.ReportSkipped {
				pass.Reportf(
					obj.Struct,
					"%s (%s) field not found in schema of `%s` object",
					field.Names[0], name, def.SchemaName,
				)
			}

			if debug {
				fmt.Printf("\t%v(%s): %v", field.Names, name, field.Type)
				fmt.Print(" -> ", "?")
				fmt.Println()
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
