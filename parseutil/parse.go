package parseutil

import (
	"go/ast"
	"reflect"
	"strings"
)

// The difference between ASCII lower 'a' and 'A'
const asciiDiff = 'a' - 'A'

func IsASCIIUpper(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func PascalToSnake(s string) string {
	if s == "" {
		return s
	}

	b := &strings.Builder{}
	b.Grow(len(s))

	sliced := []byte(s)

	for i := 0; i < len(sliced); i++ {
		c := sliced[i]

		if IsASCIIUpper(c) {
			// add dash
			// - if it's first letter
			if i != 0 {
				// - if prev is not upper
				if !IsASCIIUpper(sliced[i-1]) {
					b.WriteByte('_')
					// - or if next is not upper
				} else if i+1 < len(sliced) && !IsASCIIUpper(sliced[i+1]) {
					b.WriteByte('_')
				}
			}

			b.WriteByte(c + asciiDiff)
			continue
		}

		b.WriteByte(c)
	}

	return b.String()
}

// Parses Go tag `json:`
func ParseJSONTag(bl *ast.BasicLit) string {
	if bl == nil {
		return ""
	}

	tags := bl.Value[1 : len(bl.Value)-1] // remove quotes ``
	tag := strings.Split(reflect.StructTag(tags).Get("json"), ",")[0]

	return tag
}
