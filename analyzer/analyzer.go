package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = New()

// New creates a loglint analyzer with the provided settings.
func New() *analysis.Analyzer {
	settings := LoadSettings()

	return &analysis.Analyzer{
		Name: "loglint",
		Doc:  "checks log messages",
		Run: func(pass *analysis.Pass) (interface{}, error) {
			return run(pass, settings)
		},
	}
}

func run(pass *analysis.Pass, settings Settings) (interface{}, error) {

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
					if settings.LowercaseStart && !literalsFirstLetterIsLowercase(lits) {
						pass.Reportf(call.Pos(), "log message should start with a lowercase letter")
						return true
					}

					// rule 2: string should not contain non english letters
					if settings.EnglishOnly && literalsContainNonEnglish(lits) {
						pass.Reportf(call.Pos(), "log message must be in English")
						return true
					}

					// rule 4: string should not contain potentially sensitive data
					if settings.SensitiveData && literalsContainSensitive(lits, settings.SensitiveKeywords) {
						pass.Reportf(call.Pos(), "log message contains potentially sensitive data")
						return true
					}

					// rule 3: strign should not contain special symbols or emoji
					if settings.SpecialSymbols && literalsContainSpecial(lits) {
						pass.Reportf(call.Pos(), "log message contains special symbols or emoji")
						return true
					}
				}

				return true
			}

			// rule 1: starts with lowercase letter
			if settings.LowercaseStart && !startsWithLowercaseLetter(msg) {
				pass.Reportf(call.Pos(), "log message should start with a lowercase letter")
				return true
			}

			// rule 2: string should not contain non english letters
			if settings.EnglishOnly && isStringContainsNonEnglish(msg) {
				pass.Reportf(call.Pos(), "log message must be in English")
				return true
			}

			// rule 4: string should not contain potentially sensitive data
			if settings.SensitiveData && containsSensitive(msg, settings.SensitiveKeywords) {
				pass.Reportf(call.Pos(), "log message contains potentially sensitive data")
				return true
			}

			// rule 3: strign should not contain special symbols or emoji
			if settings.SpecialSymbols && containsSpecialSymbolsOrEmojis(msg) {
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
