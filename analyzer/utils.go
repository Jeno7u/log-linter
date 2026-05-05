package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"

	"golang.org/x/tools/go/analysis"
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

// getString extracts an unquoted string from a BasicLit if possible.
func getString(expr ast.Expr) (string, bool) {
	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		return "", false
	}

	if lit.Kind != token.STRING {
		return "", false
	}

	s, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", false
	}

	return s, true
}

// getFuncName returns the simple function name for a CallExpr.
func getFuncName(call *ast.CallExpr) string {
	switch fun := call.Fun.(type) {
	case *ast.SelectorExpr:
		return fun.Sel.Name
	case *ast.Ident:
		return fun.Name
	}
	return ""
}

// checking is this call by "log/slog"
func isSlogCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	obj := pass.TypesInfo.Uses[ident]
	if obj == nil {
		return false
	}

	pkg := obj.Pkg()
	if pkg == nil {
		return false
	}

	return pkg.Path() == "log/slog"
}

// checking is this call by "go.uber.org/zap"
func isZapCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	innerCall, ok := sel.X.(*ast.CallExpr)
	if !ok {
		return false
	}

	innerSel, ok := innerCall.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := innerSel.X.(*ast.Ident)
	if !ok {
		return false
	}

	obj := pass.TypesInfo.Uses[ident]
	if obj == nil {
		return false
	}

	pkg := obj.Pkg()
	if pkg == nil {
		return false
	}

	return pkg.Path() == "go.uber.org/zap"
}
