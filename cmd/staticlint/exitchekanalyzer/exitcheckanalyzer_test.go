package exitchekanalyzer_test

import (
	"testing"

	"github.com/koteyye/shortener/cmd/staticlint/exitchekanalyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestExitCheck(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), exitchekanalyzer.OsExitCheckAnalyzer, "./...")
}
