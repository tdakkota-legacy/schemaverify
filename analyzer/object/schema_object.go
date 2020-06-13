package object

import (
	schema "github.com/lestrrat-go/jsschema"
	"go/ast"
	"schemaverify/analyzer/parseutil"
	"schemaverify/analyzer/pragma"
)

type SchemaObjects struct {
	rewrite     map[string]string
	objects     Structs
	definitions Definitions
}

func NewSchemaObjects(root *schema.Schema) SchemaObjects {
	return SchemaObjects{
		rewrite:     map[string]string{},
		objects:     Structs{},
		definitions: root.Definitions,
	}
}

func (o SchemaObjects) InspectCallback(node ast.Node) bool {
	if v, ok := node.(*ast.GenDecl); ok {
		for _, spec := range v.Specs {
			if typeSpec, ok := spec.(*ast.TypeSpec); ok && typeSpec.Name.IsExported() { // skip unexported
				if _, ok := typeSpec.Type.(*ast.FuncType); ok { // skip functions types
					return true
				}

				pragmas := pragma.ParsePragmas(v.Doc)
				if pragmas.Skip() {
					return true
				}

				if rewrite, ok := pragmas.Rewrite(); ok {
					o.AddRewrite(typeSpec, rewrite)
				}

				o.AddObject(typeSpec)
			}
		}
	}

	return true
}

func (o SchemaObjects) AddObject(obj SchemaObject) {
	o.objects[obj.Name.Name] = obj
}

func (o SchemaObjects) AddRewrite(obj SchemaObject, schemaName string) {
	o.rewrite[obj.Name.Name] = schemaName
}

func (o SchemaObjects) FindObject(t ast.Expr) (SchemaObject, bool) {
	ident, ok := t.(*ast.Ident)
	if !ok {
		return nil, false
	}

	obj, ok := o.objects[ident.Name]
	return obj, ok
}

type DefinitionResult struct {
	SchemaName string
	*schema.Schema
}

func (o SchemaObjects) FindDefinition(structName string) (DefinitionResult, bool) {
	if rewrite, ok := o.rewrite[structName]; ok {
		sch, ok := o.definitions[rewrite]
		return DefinitionResult{SchemaName: rewrite, Schema: sch}, ok
	}

	schemaName := parseutil.PascalToSnake(structName)
	sch, ok := o.definitions[schemaName]
	return DefinitionResult{SchemaName: schemaName, Schema: sch}, ok
}

func (o SchemaObjects) ForEach(cb func(name string, obj SchemaObject) (bool, error)) error {
	for name, object := range o.objects {
		ok, err := cb(name, object)
		if err != nil {
			return err
		}

		if !ok {
			break
		}
	}

	return nil
}

func (o SchemaObjects) MapFields(list []*ast.Field) map[string]*ast.Field {
	fields := map[string]*ast.Field{}

	for _, field := range list {
		var name string

		pragmas := pragma.ParsePragmas(field.Doc)

		if rewrite, ok := pragmas.Rewrite(); ok {
			name = rewrite
		} else if v := parseutil.ParseJSONTag(field.Tag); v != "" {
			name = v
		} else if len(field.Names) > 0 {
			name = field.Names[0].Name
		} else if obj, ok := o.FindObject(field.Type); ok {
			switch def := obj.Type.(type) {
			case *ast.StructType:
				for k, v := range o.MapFields(def.Fields.List) {
					fields[k] = v
				}
			}
		}

		if name == "" {

			continue
		}

		fields[name] = field
	}

	return fields
}
