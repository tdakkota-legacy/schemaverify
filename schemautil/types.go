package schemautil

import (
	"go/types"

	schema "github.com/lestrrat-go/jsschema"
)

func GolangTypeToSchema(golangType types.Type) schema.PrimitiveType {
	v, ok := golangType.Underlying().(*types.Basic)
	if !ok {
		return schema.UnspecifiedType
	}

	switch v.Kind() {
	case types.Int,
		types.Int8,
		types.Int16,
		types.Int32,
		types.Int64,
		types.Uint,
		types.Uint8,
		types.Uint16,
		types.Uint32,
		types.Uint64:
		return schema.IntegerType

	case types.Float32, types.Float64:
		return schema.NumberType

	case types.String:
		return schema.StringType

	case types.Bool:
		return schema.BooleanType

	default:
		return schema.UnspecifiedType
	}
}

func CheckPrimitiveType(golangType types.Type, prop *schema.Schema) bool {
	casted := GolangTypeToSchema(golangType)
	return prop.Type.Contains(casted) // check that field satisfies at least one type
}
