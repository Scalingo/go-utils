package gomockgenerator

import (
	"crypto/sha1"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
)

func interfaceHash(pkg, iName string) (string, error) {
	sig, err := interfaceSignature(pkg, iName)
	if err != nil {
		return "", errors.Wrapf(err, "fail to get interface signature for %s", iName)
	}
	hash := sha1.Sum([]byte(sig))
	return fmt.Sprintf("% x", hash), nil
}

func interfaceSignature(pkg, iName string) (string, error) {
	fullPath := path.Join(os.Getenv("GOPATH"), "src", pkg)
	fileSet := token.NewFileSet()
	packages, err := parser.ParseDir(fileSet, fullPath, func(info os.FileInfo) bool {
		return !strings.HasSuffix(info.Name(), "_test.go")
	}, 0)

	if err != nil {
		return "", errors.Wrap(err, "fail to parse package")
	}

	if len(packages) == 0 {
		return "", errors.Errorf("no package found in %v for interface %v", fullPath, iName)
	}

	if len(packages) != 1 {
		for name := range packages {
			fmt.Println(name)
		}
		return "", errors.New("too many packages")
	}

	for _, pck := range packages {
		for _, f := range pck.Files {
			for _, decl := range f.Decls {
				if gen, ok := decl.(*ast.GenDecl); ok {
					if gen.Tok == token.TYPE {
						for _, spec := range gen.Specs {
							if typeSpec, ok := spec.(*ast.TypeSpec); ok {
								if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
									if typeSpec.Name.String() == iName {
										var interfaceSig string
										for _, m := range interfaceType.Methods.List {
											switch v := m.Type.(type) {
											case *ast.Ident:
												interfaceSig = fmt.Sprintf("%s\nInterface{%s}\n", interfaceSig, v.Name)
											case *ast.FuncType:
												methodName := m.Names[0].String()
												var methodType string
												if v.Results != nil {
													for _, res := range v.Results.List {
														methodType = fmt.Sprintf("%s,%s", methodType, fieldToString(res.Type))
													}
												}

												var methodParams string
												if v.Params != nil {
													for _, param := range v.Params.List {
														methodParams = fmt.Sprintf("%s,%s", methodParams, fieldToString(param.Type))
													}
												}
												interfaceSig = fmt.Sprintf("%s\n%s(%s)(%s)\n", interfaceSig, methodName, methodParams, methodType)

											default:
												panic(fmt.Sprintf("Unexpected AST type: %T", v))
											}
										}
										return interfaceSig, nil
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return "", errors.New("not found")
}

func fieldToString(field ast.Expr) string {
	if starExpr, ok := field.(*ast.StarExpr); ok {
		return "*" + fieldToString(starExpr.X)
	}

	if mapExpr, ok := field.(*ast.MapType); ok {
		return fmt.Sprintf("map[%s]%s", fieldToString(mapExpr.Key), fieldToString(mapExpr.Value))
	}

	if arrayExpr, ok := field.(*ast.ArrayType); ok {
		return fmt.Sprintf("[]%v", fieldToString(arrayExpr.Elt))
	}

	if ellipsisExpr, ok := field.(*ast.Ellipsis); ok {
		return fmt.Sprintf("...%v", fieldToString(ellipsisExpr.Elt))
	}

	if _, ok := field.(*ast.InterfaceType); ok {
		// TODO: It's the only case I can think of, but they might be others.
		return "interface{}"
	}

	return fmt.Sprintf("%v", field)
}
