package analyzer

import (
	"go/ast"

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
				// not a basic literal: extract literal operands from concatenation
				lits := collectStringLiterals(call.Args[0])
				if len(lits) > 0 {
					// rule 1: starts with lowercase letter
					if !literalsFirstLetterIsLowercase(lits) {
						pass.Reportf(call.Pos(), "log message should start with a lowercase letter")
						return true
					}

					// rule 2: string should not contain non english letters
					if literalsContainNonEnglish(lits) {
						pass.Reportf(call.Pos(), "log message must be in English")
						return true
					}

					// rule 4: string should not contain potentially sensitive data
					if literalsContainSensitive(lits) {
						pass.Reportf(call.Pos(), "log message contains potentially sensitive data")
						return true
					}

					// rule 3: strign should not contain special symbols or emoji
					if literalsContainSpecial(lits) {
						pass.Reportf(call.Pos(), "log message contains special symbols or emoji")
						return true
					}
				}

				return true
			}

			// rule 1: starts with lowercase letter
			if !startsWithLowercaseLetter(msg) {
				pass.Reportf(call.Pos(), "log message should start with a lowercase letter")
				return true
			}

			// rule 2: string should not contain non english letters
			if isStringContainsNonEnglish(msg) {
				pass.Reportf(call.Pos(), "log message must be in English")
				return true
			}

			// rule 4: string should not contain potentially sensitive data
			if containsSensitive(msg) {
				pass.Reportf(call.Pos(), "log message contains potentially sensitive data")
				return true
			}

			// rule 3: strign should not contain special symbols or emoji
			if containsSpecialSymbolsOrEmojis(msg) {
				pass.Reportf(call.Pos(), "log message contains special symbols or emoji")
				return true
			}

			return true
		})
	}

	return nil, nil
}

// checks is call was by slog or zap logger
func isLoggerCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	return isSlogCall(pass, call) || isZapCall(pass, call)
}
