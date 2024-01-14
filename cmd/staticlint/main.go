package main

import (
	errcheckanalyzer "github.com/koteyye/shortener/cmd/staticlint/errcheckanalyzer"
	"github.com/koteyye/shortener/cmd/staticlint/exitchekanalyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

var excludeStyleChecks = map[string]struct{}{
	"ST1000": {},
	"ST1020": {},
	"ST1021": {},
	"ST1022": {},
}

func main() {
	analyzers := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		assign.Analyzer,
		errcheckanalyzer.ErrCheckAnalyzer,
		exitchekanalyzer.OsExitCheckAnalyzer,
		structtag.Analyzer,
	}

	for _, v := range staticcheck.Analyzers {
		if _, ok := excludeStyleChecks[v.Analyzer.Name]; !ok {
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	multichecker.Main(
		analyzers...,
	)
}
