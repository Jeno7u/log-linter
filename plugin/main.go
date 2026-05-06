package plugin

import (
	"github.com/Jeno7u/log-linter/analyzer"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("loglint", New)
}

// creates the loglint plugin instance
func New(conf any) (register.LinterPlugin, error) {
	return &plugin{analyzer: analyzer.New(analyzer.LoadSettings())}, nil
}

// plugin exposes the analyzer to golangci-lint
type plugin struct {
	analyzer *analysis.Analyzer
}

// returns the analyzers provided by this plugin
func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{p.analyzer}, nil
}

// reports the load mode required by the analyzer
func (p *plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
