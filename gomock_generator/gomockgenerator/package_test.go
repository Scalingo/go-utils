package gomockgenerator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
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
		"channel":    {expression: "chan Event", want: "chan Event"},
		"send channel": {
			expression: "chan<- *Event",
			want:       "chan<- *Event",
		},
		"receive channel": {
			expression: "<-chan *Event",
			want:       "<-chan *Event",
		},
		"function": {
			expression: "func(value string, values ...Event) (int, error)",
			want:       "func(string,...Event) (int,error)",
		},
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

		funcDecl, ok := file.Decls[0].(*ast.FuncDecl)
		if !ok {
			t.Errorf("parsed declaration for %q has type %T, want *ast.FuncDecl", expression, file.Decls[0])
			return nil
		}
		return funcDecl.Type.Params.List[0].Type
	}

	expr, err := parser.ParseExpr(expression)
	if err != nil {
		t.Errorf("parse expression %q: %v", expression, err)
		return nil
	}

	return expr
}

func TestInterfaceSignature(t *testing.T) {
	t.Run("resolves a standard library package", func(t *testing.T) {
		signature, err := interfaceSignature(t.Context(), "net", "Conn")
		if err != nil {
			t.Fatalf("interfaceSignature returned error: %v", err)
		}

		if signature == "" {
			t.Fatal("interfaceSignature returned an empty signature")
		}

		for _, method := range []string{"Read", "Write", "Close", "SetDeadline"} {
			if !strings.Contains(signature, method) {
				t.Fatalf("signature %q does not contain method %q", signature, method)
			}
		}
	})

	t.Run("serializes methods and embedded interfaces", func(t *testing.T) {
		pkg := fixturePackage("service")

		signature, err := interfaceSignature(t.Context(), pkg, "Service")
		if err != nil {
			t.Fatalf("interfaceSignature returned error: %v", err)
		}

		want := "\nInterface{Embedded}\n\nPing()()\n\nTransform(,string,[]byte)(,int,error)\n"
		if signature != want {
			t.Errorf("interfaceSignature() = %q, want %q", signature, want)
		}
	})

	t.Run("forces regeneration for a qualified embedded interface", func(t *testing.T) {
		pkg := fixturePackage("qualified")

		signature, err := interfaceSignature(t.Context(), pkg, "Service")
		if err != nil {
			t.Fatalf("interfaceSignature returned error: %v", err)
		}
		if signature != "FORCE_REGENERATE" {
			t.Errorf("interfaceSignature() = %q, want %q", signature, "FORCE_REGENERATE")
		}
	})

	t.Run("reports a missing interface", func(t *testing.T) {
		pkg := fixturePackage("missing")

		_, err := interfaceSignature(t.Context(), pkg, "Missing")
		if err == nil {
			t.Fatal("interfaceSignature returned no error")
		}
		if !strings.Contains(err.Error(), "interface Missing not found") {
			t.Errorf("interfaceSignature error = %q, want missing-interface context", err)
		}
	})

	t.Run("changes when method types change", func(t *testing.T) {
		stringPackage := fixturePackage("string_parameter")
		intPackage := fixturePackage("int_parameter")

		stringSignature, err := interfaceSignature(t.Context(), stringPackage, "Service")
		if err != nil {
			t.Fatalf("get string signature: %v", err)
		}
		intSignature, err := interfaceSignature(t.Context(), intPackage, "Service")
		if err != nil {
			t.Fatalf("get int signature: %v", err)
		}

		if stringSignature == intSignature {
			t.Fatalf("method type change produced the same signature %q", stringSignature)
		}
	})

	t.Run("serializes channel and function types consistently", func(t *testing.T) {
		signature, err := interfaceSignature(t.Context(), fixturePackage("channels_and_functions"), "Service")
		if err != nil {
			t.Fatalf("interfaceSignature returned error: %v", err)
		}

		want := "\nEvents()(,<-chan *Event)\n\nSubscribe(,func(&{context Context},*&{time Timer},string) error)(,chan<- Event,func() error)\n"
		if signature != want {
			t.Errorf("interfaceSignature() = %q, want %q", signature, want)
		}
	})
}

func TestInterfaceHash(t *testing.T) {
	t.Run("hashes the interface signature", func(t *testing.T) {
		hash, err := interfaceHash(t.Context(), fixturePackage("service"), "Service")
		if err != nil {
			t.Fatalf("interfaceHash returned error: %v", err)
		}

		const want = "fd 79 e0 1f 61 c6 be 27 5c 24 01 19 b5 4f 79 8d a1 ac ac 10"
		if hash != want {
			t.Errorf("interfaceHash() = %q, want %q", hash, want)
		}
	})

	t.Run("preserves the force regeneration marker", func(t *testing.T) {
		hash, err := interfaceHash(t.Context(), fixturePackage("qualified"), "Service")
		if err != nil {
			t.Fatalf("interfaceHash returned error: %v", err)
		}

		if hash != "FORCE_REGENERATE" {
			t.Errorf("interfaceHash() = %q, want %q", hash, "FORCE_REGENERATE")
		}
	})

	t.Run("hashes channel and function types consistently", func(t *testing.T) {
		hash, err := interfaceHash(t.Context(), fixturePackage("channels_and_functions"), "Service")
		if err != nil {
			t.Fatalf("interfaceHash returned error: %v", err)
		}

		const want = "15 4e e3 29 40 28 ae 6e d5 54 da e7 5e 59 97 f7 7e 97 69 98"
		if hash != want {
			t.Errorf("interfaceHash() = %q, want %q", hash, want)
		}
	})
}

func fixturePackage(name string) string {
	return filepath.Join("fixtures", name)
}
