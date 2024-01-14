package errcheckanalyzer_test

import (
	"testing"

	"github.com/koteyye/shortener/cmd/staticlint/errcheckanalyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestErrCheckAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), errcheckanalyzer.ErrCheckAnalyzer, "./...")
}
