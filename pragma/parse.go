package pragma

import (
	"go/ast"
	"strings"
)

const linterPragmaPrefix = "//schemaverify:"

type Parser struct {
	prefix       string
	nolintPrefix []string
}

func NewParser(prefix string, nolintPrefix []string) Parser {
	return Parser{prefix: prefix, nolintPrefix: nolintPrefix}
}

func (p Parser) isNoLintPrefix(s string) bool {
	for i := range p.nolintPrefix {
		if strings.HasPrefix(s, p.nolintPrefix[i]) {
			return true
		}
	}

	return false
}

func (p Parser) isPragmaComment(s string) bool {
	return strings.HasPrefix(s, p.prefix)
}

type Comment struct {
	Value  string
	NoLint bool
}

func (p Parser) ParseCommentGroup(doc *ast.CommentGroup) []Comment {
	if doc == nil {
		return nil
	}

	var result []Comment
	for _, comment := range doc.List {
		switch {
		case p.isPragmaComment(comment.Text):
			result = append(result, Comment{
				Value: strings.TrimSpace(strings.TrimPrefix(comment.Text, linterPragmaPrefix)),
			})

		case p.isNoLintPrefix(comment.Text):
			result = append(result, Comment{
				Value:  "",
				NoLint: true,
			})
		}
	}

	return result
}

func parsePragma(pair string) (string, string, bool) {
	if pair == "" {
		return "", "", false
	}

	if strings.Index(pair, "=") < 0 {
		return pair, "", true
	}

	split := strings.SplitN(pair, "=", 2)
	if len(split) == 2 {
		return split[0], split[1], true
	}

	return "", "", false
}

func (p Parser) ParsePragmas(comments []Comment) Pragmas {
	result := Pragmas{}

	for _, pragma := range comments {
		if pragma.NoLint {
			result[SkipPragma] = ""
			continue
		}

		pairs := strings.SplitN(pragma.Value, ",", 1)

		for _, pair := range pairs {
			if k, v, ok := parsePragma(pair); ok {
				result[k] = v
			}
		}
	}

	return result
}

func ParsePragmas(doc *ast.CommentGroup) Pragmas {
	parser := NewParser(linterPragmaPrefix, []string{
		"// BUG",
	})

	comments := parser.ParseCommentGroup(doc)
	return parser.ParsePragmas(comments)
}
