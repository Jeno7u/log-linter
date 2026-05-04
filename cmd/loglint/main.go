package main

import (
	"github.com/Jeno7u/log-linter/analyzer"

	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		analyzer.Analyzer,
	)
}
