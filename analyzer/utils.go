package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
)

var sensitiveKeywords = []string{
	"password", "token", "api_key", "apikey", "secret",
}

// returns true if given rune is english letter
func isEnglishLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

// returns true if given rune is number
func isNumber(r rune) bool {
	return r >= '0' && r <= '9'
}

// collectStringLiterals extracts string literal operands from a concatenation expression
// (left-to-right). It returns a slice of unquoted literal strings.
func collectStringLiterals(expr ast.Expr) []string {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			s, err := strconv.Unquote(e.Value)
			if err == nil {
				return []string{s}
			}
		}
		return nil
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			left := collectStringLiterals(e.X)
			right := collectStringLiterals(e.Y)
			return append(left, right...)
		}
		return nil
	case *ast.ParenExpr:
		return collectStringLiterals(e.X)
	default:
		return nil
	}
}
