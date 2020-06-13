package object

import "go/ast"

type SchemaObject *ast.TypeSpec

type Structs map[string]SchemaObject

func (s Structs) Find(name string) (SchemaObject, bool) {
	def, ok := s[name]
	return def, ok
}
