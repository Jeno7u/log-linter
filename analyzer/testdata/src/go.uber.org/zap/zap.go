package zap

// Minimal fake zap package for analysistest. Only provides functions used in testdata.

type logger struct{}

func (l *logger) Info(msg string)  {}
func (l *logger) Error(msg string) {}
func (l *logger) Warn(msg string)  {}
func (l *logger) Debug(msg string) {}

func L() *logger { return &logger{} }
