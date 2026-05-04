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

			// фильтруем только логгеры
			if !isLoggerCall(funcName) {
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

func isLoggerCall(name string) bool {
	switch name {
	case "Info", "Error", "Warn", "Debug":
		return true
	}
	return false
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
