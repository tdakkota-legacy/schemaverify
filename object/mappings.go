package object

import (
	"go/ast"

	"github.com/tdakkota/schemaverify/parseutil"
	"github.com/tdakkota/schemaverify/pragma"

	schema "github.com/lestrrat-go/jsschema"
)

type Mappings struct {
	structMapping     map[string]int
	schemaMapping     map[string]int
	schemaDefinitions Definitions
	pairs             []Pair
}

func NewMapping(root *schema.Schema) Mappings {
	return Mappings{
		structMapping:     map[string]int{},
		schemaMapping:     map[string]int{},
		schemaDefinitions: root.Definitions,
	}
}

func (o *Mappings) InspectCallback(node ast.Node) bool {
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

				var schemaName string
				if rewrite, ok := pragmas.Rewrite(); ok {
					schemaName = rewrite
				} else {
					schemaName = parseutil.PascalToSnake(typeSpec.Name.Name)
				}

				def, ok := o.schemaDefinitions.Find(schemaName)
				if !ok {
					def = nil
				}

				pair := Pair{
					Object: typeSpec,
					Definition: Definition{
						SchemaName: schemaName,
						Schema:     def,
					},
				}

				o.AddPair(pair)
			}
		}
	}

	return true
}

func (o *Mappings) AddPair(pair Pair) {
	o.pairs = append(o.pairs, pair)

	index := len(o.pairs) - 1
	o.structMapping[pair.Object.Name.Name] = index
	o.schemaMapping[pair.Definition.SchemaName] = index
}

func (o Mappings) FindByExpr(t ast.Expr) (Pair, bool) {
	ident, ok := t.(*ast.Ident)
	if !ok {
		return Pair{}, false
	}

	return o.FindByStructName(ident.Name)
}

func (o Mappings) FindByStructName(name string) (Pair, bool) {
	return o.findByMapping(o.structMapping, name)
}

func (o Mappings) FindBySchemaName(name string) (Pair, bool) {
	return o.findByMapping(o.schemaMapping, name)
}

func (o Mappings) findByMapping(mapping map[string]int, name string) (Pair, bool) {
	index, ok := mapping[name]
	if !ok {
		return Pair{}, false
	}

	return o.pairs[index], true
}

func (o Mappings) ForEach(cb func(index int, pair Pair) (bool, error)) error {
	for i, pair := range o.pairs {
		ok, err := cb(i, pair)
		if err != nil {
			return err
		}

		if !ok {
			break
		}
	}

	return nil
}

func (o Mappings) MapFields(list []*ast.Field) map[string]*ast.Field {
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
		} else if pair, ok := o.FindByExpr(field.Type); ok {
			switch def := pair.Object.Type.(type) {
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
