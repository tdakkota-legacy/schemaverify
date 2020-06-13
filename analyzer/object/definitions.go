package object

import (
	"schemaverify/analyzer/parseutil"

	schema "github.com/lestrrat-go/jsschema"
)

type Definitions map[string]*schema.Schema

func (defs Definitions) Find(structName string, alias []string) (*schema.Schema, bool) {
	names := append([]string{parseutil.PascalToSnake(structName)}, alias...)

	for _, name := range names {
		def, ok := defs[name]
		if ok {
			return def, true
		}
	}

	return nil, false
}

func (defs Definitions) FindAll(structName string, alias []string) (result []*schema.Schema) {
	names := append([]string{parseutil.PascalToSnake(structName)}, alias...)

	for _, name := range names {
		def, ok := defs[name]
		if ok {
			result = append(result, def)
		}
	}

	return result
}
