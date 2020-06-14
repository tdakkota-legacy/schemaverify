package object

import (
	schema "github.com/lestrrat-go/jsschema"
)

type Properties struct {
	Props    map[string]*schema.Schema
	Required []string
}

func MapProperties(sch *schema.Schema) (Properties, error) {
	props := make(map[string]*schema.Schema, len(sch.Properties))
	required := make([]string, len(sch.Required))

	for k, v := range sch.Properties {
		props[k] = v
	}
	copy(required, sch.Required)

	for i := range sch.AllOf {
		resolved, err := sch.AllOf[i].Resolve(nil)
		if err != nil {
			return Properties{}, err
		}

		child, err := MapProperties(resolved)
		if err != nil {
			return Properties{}, err
		}

		for k, v := range child.Props {
			props[k] = v
		}
		required = append(required, child.Required...)
	}

	return Properties{
		Props:    props,
		Required: required,
	}, nil
}
