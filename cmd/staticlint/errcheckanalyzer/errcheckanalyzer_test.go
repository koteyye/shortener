package errcheckanalyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/koteyye/shortener/cmd/staticlint/errcheckanalyzer"
)

func TestErrCheckAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), errcheckanalyzer.ErrCheckAnalyzer, "./...")
}
