package object

import (
	"fmt"
	"go/ast"
	"go/types"
	"path"
	"schemaverify/analyzer/schemautil"

	schema "github.com/lestrrat-go/jsschema"
)

func (o Verifier) verifyField(
	types *types.Info,
	field *ast.Field,
	prop *schema.Schema,
) (string, error) {
	return o.verifyType(types, field.Type, prop)
}

func (o Verifier) verifyType(
	types *types.Info,
	typ ast.Expr,
	prop *schema.Schema,
) (string, error) {
	if prop.Type.Len() != 0 {
		return o.verifyPrimitiveType(types, typ, prop)
	}

	return o.verifyReference(types, typ, prop)
}

func (o Verifier) verifyReference(
	types *types.Info,
	typ ast.Expr,
	prop *schema.Schema,
) (string, error) {
	_, object := path.Split(prop.Reference)

	pair, ok := o.Objects.FindBySchemaName(object)
	if !ok {
		return "", fmt.Errorf("struct for \"%s\" reference not found", prop.Reference)
	}

	ident, ok := typ.(*ast.Ident)
	if !ok {
		return "", fmt.Errorf("struct for \"%s\" reference not found", prop.Reference)
	}

	if pair.Object.Name.Name != ident.Name {
		return "", fmt.Errorf("expected %s, got %s", pair.Object.Name.Name, ident.Name)
	}

	return prop.Reference, nil
}

func (o Verifier) verifyPrimitiveType(
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

		schemaTypeName, err := o.verifyType(types, sliceType.Elt, prop.Items.Schemas[0])
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
