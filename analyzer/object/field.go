package object

import (
	"fmt"
	"go/ast"
	"go/types"
	"schemaverify/analyzer/schemautil"

	schema "github.com/lestrrat-go/jsschema"
)

func (o Verifier) verifyField(
	types *types.Info,
	typ ast.Expr,
	prop *schema.Schema,
) (string, error) {
	if prop.Type.Len() != 0 {
		return o.verifySchemaType(types, typ, prop)
	}

	return prop.Reference, nil
}

func (o Verifier) verifySchemaType(
	types *types.Info,
	typ ast.Expr,
	prop *schema.Schema,
) (string, error) {
	if _, ok := typ.(*ast.InterfaceType); ok {
		return "interface{}", nil
	}

	prop, err := prop.Resolve(nil)
	if err != nil {
		return "", err
	}

	if prop.Type.Contains(schema.ArrayType) {
		sliceType, ok := typ.(*ast.ArrayType)
		if !ok {
			return "", fmt.Errorf("expected array type instead of %v", typ)
		}

		schemaTypeName, err := o.verifyField(types, sliceType.Elt, prop.Items.Schemas[0])
		if err != nil {
			return "", err
		}
		return "[]" + schemaTypeName, nil
	}

	if !schemautil.CheckPrimitiveType(types.TypeOf(typ), prop) {
		var expectedType string
		if len(prop.Type) == 1 {
			expectedType = prop.Type[0].String()
		} else {
			expectedType = fmt.Sprint(prop.Type)
		}

		return "", fmt.Errorf("expected %s instead of %v", expectedType, typ)
	}

	return prop.Type[0].String(), nil
}
