package object

import (
	schema "github.com/lestrrat-go/jsschema"
)

type Definitions map[string]*schema.Schema

func (defs Definitions) Find(name string) (*schema.Schema, bool) {
	def, ok := defs[name]
	return def, ok
}

func (defs Definitions) FindAll(names ...string) (result []*schema.Schema) {
	for _, name := range names {
		def, ok := defs[name]
		if ok {
			result = append(result, def)
		}
	}

	return result
}
