package exitchekanalyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// OsExitCheckAnalyzer прямое использование os.Exit
// в функции main
var OsExitCheckAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check the use os.exit in main func",
	Run:  runExitCheck,
}

func runExitCheck(pass *analysis.Pass) (any, error) {
	expr := func(x *ast.SelectorExpr) bool {
		pkg, ok := x.X.(*ast.Ident)
		if ok && pkg.Name == "os" && x.Sel.Name == "Exit" {
			pass.Reportf(x.Pos(), "calling os.Exit in main")
			return false
		}
		return true
	}

	for _, file := range pass.Files {
		if file.Name.Name != "main" {
			continue
		}
		ast.Inspect(file, func(n ast.Node) bool {
			if x, ok := n.(*ast.SelectorExpr); ok {
				return expr(x)
			}
			return true
		})
	}
	return nil, nil
}
