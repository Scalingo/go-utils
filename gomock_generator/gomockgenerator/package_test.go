package gomockgenerator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestFieldToString(t *testing.T) {
	tests := map[string]struct {
		expression string
		want       string
	}{
		"identifier": {expression: "Event", want: "Event"},
		"pointer":    {expression: "*Event", want: "*Event"},
		"map":        {expression: "map[string]*Event", want: "map[string]*Event"},
		"slice":      {expression: "[]Event", want: "[]Event"},
		"ellipsis":   {expression: "...Event", want: "...Event"},
		"interface":  {expression: "interface{}", want: "interface{}"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			firstExpr := parseFieldExpression(t, tt.expression)
			if firstExpr == nil {
				return
			}
			secondExpr := parseFieldExpression(t, tt.expression)
			if secondExpr == nil {
				return
			}

			first := fieldToString(firstExpr)
			second := fieldToString(secondExpr)
			if first != second {
				t.Errorf("fieldToString is unstable: first %q, second %q", first, second)
				return
			}
			if first != tt.want {
				t.Errorf("fieldToString() = %q, want %q", first, tt.want)
			}
		})
	}
}

func parseFieldExpression(t *testing.T, expression string) ast.Expr {
	t.Helper()

	if strings.HasPrefix(expression, "...") {
		file, err := parser.ParseFile(token.NewFileSet(), "", "package p\nfunc f(value "+expression+") {}", 0)
		if err != nil {
			t.Errorf("parse parameter expression %q: %v", expression, err)
			return nil
		}

		funcDecl := file.Decls[0].(*ast.FuncDecl)
		return funcDecl.Type.Params.List[0].Type
	}

	expr, err := parser.ParseExpr(expression)
	if err != nil {
		t.Errorf("parse expression %q: %v", expression, err)
		return nil
	}

	return expr
}

func TestInterfaceSignatureResolvesStandardLibraryPackage(t *testing.T) {
	sig, err := interfaceSignature(t.Context(), "net", "Conn")
	if err != nil {
		t.Fatalf("interfaceSignature returned error: %v", err)
	}

	if sig == "" {
		t.Fatal("interfaceSignature returned an empty signature")
	}

	for _, method := range []string{"Read", "Write", "Close", "SetDeadline"} {
		if !strings.Contains(sig, method) {
			t.Fatalf("signature %q does not contain method %q", sig, method)
		}
	}
}
