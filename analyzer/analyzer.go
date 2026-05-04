package analyzer

import (
	"go/ast"
	"go/token"

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

			// получаем имя вызываемой функции
			funcName := getFuncName(call)

			// фильтруем только нужные нам логгеры
			if !isLoggerCall(pass, call) {
				return true
			}

			// достаём первый аргумент
			if len(call.Args) == 0 {
				return true
			}

			msg, ok := getString(call.Args[0])
			if !ok {
				return true
			}

			pass.Reportf(call.Pos(), "log message: %s", msg)

			return true
		})
	}

	return nil, nil
}

func getFuncName(call *ast.CallExpr) string {
	switch fun := call.Fun.(type) {

	case *ast.SelectorExpr:
		// slog.Info или logger.Info
		return fun.Sel.Name

	case *ast.Ident:
		// Println
		return fun.Name
	}

	return ""
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

	return lit.Value, true
}

// проверяем, вызов ли это "log/slog"
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

// проверяем, вызов ли это "go.uber.org/zap"
func isZapCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	// zap.L().Info → sel.X это CallExpr
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
