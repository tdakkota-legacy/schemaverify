package object

import (
	schema "github.com/lestrrat-go/jsschema"
)

type Definitions map[string]*schema.Schema

func (defs Definitions) Find(name string) (*schema.Schema, bool) {
	def, ok := defs[name]
	return def, ok
}
