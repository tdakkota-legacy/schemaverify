package pragma

const (
	SchemaNamePragma = "schema_name"
	SkipPragma       = "skip"
)

type Pragmas map[string]string

func (p Pragmas) Skip() bool {
	_, ok := p[SkipPragma]
	return ok
}

func (p Pragmas) Rewrite() (string, bool) {
	v, ok := p[SchemaNamePragma]

	return v, ok
}
