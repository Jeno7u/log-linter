package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "loglint",
	Doc:  "checks log messages",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {

		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			// checking that call is by logger
			if !isLoggerCall(pass, call) {
				return true
			}

			// taking first argument
			if len(call.Args) == 0 {
				return true
			}

			msg, ok := getString(call.Args[0])
			if !ok {
				return true
			}

			// rule 1: should start with lowercase letter
			if !isStringNonEmpty(msg) || !firstLetterIsLowercase(msg) {
				pass.Reportf(call.Pos(), "log message should start with a lowercase letter: %s", msg)
			}

			// rule 2: only english
			if isStringContainsNonEnglish(msg) {
				pass.Reportf(call.Pos(), "log message contains letters that are non english: %s", msg)
			}

			// rule 3: no special symbols or emojis
			if containsSpecialSymbolsOrEmojis(msg) {
				pass.Reportf(call.Pos(), "log message contains special symbols or emoji: %s", msg)
			}

			// rule 4: no sensitive data
			if containsSensitive(msg) {
				pass.Reportf(call.Pos(), "log message contains potentially sensitive data: %s", msg)
			}

			return true
		})
	}

	return nil, nil
}

func isLoggerCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	return isSlogCall(pass, call) || isZapCall(pass, call)
}

func getString(expr ast.Expr) (string, bool) {
	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		return "", false
	}

	if lit.Kind != token.STRING {
		return "", false
	}

	// lit.Value includes surrounding quotes (`"..."` or ` `...` `).
	s, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", false
	}

	return s, true
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
