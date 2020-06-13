package pragma

import "strings"

const (
	SchemaNamePragma  = "schema_name"
	SchemaNamesPragma = "schema_names"
	SkipPragma        = "skip"
)

type Pragmas map[string]string

func (p Pragmas) Skip() bool {
	_, ok := p[SkipPragma]
	return ok
}

func (p Pragmas) Alias() (string, bool) {
	v, ok := p[SchemaNamePragma]

	return v, ok
}

func (p Pragmas) Aliases() []string {
	var alias []string
	if v, ok := p.Alias(); ok {
		alias = append(alias, v)
	}

	if v, ok := p[SchemaNamesPragma]; ok {
		alias = append(alias, strings.Split(v, ",")...)
	}

	return alias
}
