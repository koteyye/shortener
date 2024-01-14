package errcheckanalyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// ErrCheckAnalyzer обнаруживает необработанные ошибки в коде,
// ошибка считается необработанной, если присваивается как '_'.
var ErrCheckAnalyzer = &analysis.Analyzer{
	Name: "errcheck",
	Doc:  "check for unchecked errors",
	Run:  runErrCheck,
}

func runErrCheck(pass *analysis.Pass) (any, error) {
	expr := func(x *ast.ExprStmt) {
		if call, ok := x.X.(*ast.CallExpr); ok {
			if isReturnError(pass, call) {
				pass.Reportf(x.Pos(), "expression returns unchecked error")
			}
		}
	}
	tuplefunc := func(x *ast.AssignStmt) {
		if call, ok := x.Rhs[0].(*ast.CallExpr); ok {
			results := resultErrors(pass, call)
			for i := 0; i < len(x.Lhs); i++ {
				if id, ok := x.Lhs[i].(*ast.Ident); ok && id.Name == "_" && results[i] {
					pass.Reportf(id.NamePos, "assignment with unchecked error")
				}
			}
		}
	}
	errFunc := func(x *ast.AssignStmt) {
		for i := 0; i < len(x.Lhs); i++ {
			if id, ok := x.Lhs[i].(*ast.Ident); ok {
				if call, ok := x.Rhs[i].(*ast.CallExpr); ok {
					if id.Name == "_" && isReturnError(pass, call) {
						pass.Reportf(id.NamePos, "assignment with unchecked error")
					}
				}
			}
		}
	}
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.ExprStmt:
				expr(x)
			case *ast.GoStmt:
				if isReturnError(pass, x.Call) {
					pass.Reportf(x.Pos(), "go statement with unchecked error")
				}
			case *ast.DeferStmt:
				if isReturnError(pass, x.Call) {
					pass.Reportf(x.Pos(), "defer with unchecked error")
				}
			case *ast.AssignStmt:
				if len(x.Rhs) == 1 {
					tuplefunc(x)
				} else {
					errFunc(x)
				}
			}
			return true
		})
	}
	return nil, nil
}

var errorType = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

func isErrorType(t types.Type) bool {
	return types.Implements(t, errorType)
}

func resultErrors(pass *analysis.Pass, call *ast.CallExpr) []bool {
	switch t := pass.TypesInfo.Types[call].Type.(type) {
	case *types.Named:
		return []bool{isErrorType(t)}
	case *types.Pointer:
		return []bool{isErrorType(t)}
	case *types.Tuple:
		s := make([]bool, t.Len())
		for i := 0; i < t.Len(); i++ {
			switch mt := t.At(i).Type().(type) {
			case *types.Named:
				s[i] = isErrorType(mt)
			case *types.Pointer:
				s[i] = isErrorType(mt)
			}
		}
		return s
	}
	return []bool{false}
}

func isReturnError(pass *analysis.Pass, call *ast.CallExpr) bool {
	for _, isError := range resultErrors(pass, call) {
		if isError {
			return true
		}
	}
	return false
}
