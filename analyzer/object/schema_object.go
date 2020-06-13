package object

import (
	"go/ast"
	"schemaverify/analyzer/parseutil"
	"schemaverify/analyzer/pragma"
)

type SchemaObject *ast.TypeSpec

type SchemaObjects struct {
	alias   map[string][]string
	objects map[string]SchemaObject
}

func NewSchemaObjects() SchemaObjects {
	return SchemaObjects{
		alias:   map[string][]string{},
		objects: map[string]SchemaObject{},
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

				o.AddObject(typeSpec, pragmas.Aliases()...)
			}
		}
	}

	return true
}

func (o SchemaObjects) AddObject(obj SchemaObject, aliases ...string) {
	o.objects[obj.Name.Name] = obj
	o.alias[obj.Name.Name] = aliases
}

func (o SchemaObjects) FindObject(t ast.Expr) (SchemaObject, bool) {
	ident, ok := t.(*ast.Ident)
	if !ok {
		return nil, false
	}

	obj, ok := o.objects[ident.Name]
	return obj, ok
}

func (o SchemaObjects) ForEach(cb func(name string, alias []string, obj SchemaObject) (bool, error)) error {
	for name, object := range o.objects {
		ok, err := cb(name, o.alias[name], object)
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

		if v := parseutil.ParseJsonTag(field.Tag); v != "" {
			name = v
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
