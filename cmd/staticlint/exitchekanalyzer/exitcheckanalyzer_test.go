package exitchekanalyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/koteyye/shortener/cmd/staticlint/exitchekanalyzer"
)

func TestExitCheck(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), exitchekanalyzer.OsExitCheckAnalyzer, "./...")
}
