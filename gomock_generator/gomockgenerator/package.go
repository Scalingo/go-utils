package gomockgenerator

import (
	"context"
	"crypto/sha1"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Scalingo/go-utils/errors/v3"
)

func interfaceHash(ctx context.Context, pkg, iName string) (string, error) {
	sig, err := interfaceSignature(ctx, pkg, iName)
	if err != nil {
		return "", errors.Wrapf(ctx, err, "get interface signature for %s:%s", pkg, iName)
	}
	if sig == "FORCE_REGENERATE" {
		return sig, nil
	}
	hash := sha1.Sum([]byte(sig))
	return fmt.Sprintf("% x", hash), nil
}

func interfaceSignature(ctx context.Context, pkg, iName string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(ctx, err, "get current working directory")
	}

	// Resolve packages in this order:
	// 1. the project tree relative to the current working directory,
	// 2. the local vendor directory,
	// 3. the package location reported by `go list` for standard library and module packages.
	fullPath := pkg
	if !filepath.IsAbs(fullPath) {
		localPkg := filepath.Join(cwd, fullPath)
		vendoredPkg := filepath.Join(cwd, "vendor", fullPath)
		switch {
		case isDir(localPkg):
			fullPath = localPkg
		case isDir(vendoredPkg):
			fullPath = vendoredPkg
		default:
			resolvedPkg, err := resolvePackageDir(ctx, cwd, pkg)
			if err != nil {
				// Fall back to the local path when `go list` cannot resolve the package.
				fullPath = localPkg
			} else {
				fullPath = resolvedPkg
			}
		}
	}

	fileSet := token.NewFileSet()
	packages, err := parser.ParseDir(fileSet, fullPath, func(info os.FileInfo) bool {
		return !strings.HasSuffix(info.Name(), "_test.go")
	}, 0)

	if err != nil {
		return "", errors.Wrap(ctx, err, "parse package")
	}

	if len(packages) == 0 {
		return "", errors.Newf(ctx, "no package found in %v for interface %v", fullPath, iName)
	}

	if len(packages) != 1 {
		for name := range packages {
			fmt.Println(name)
		}
		return "", errors.New(ctx, "too many packages")
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
											case *ast.SelectorExpr:
												// If there is a selector expr (if the interface calls other interfaces: force a regeneration)
												// Implementing a real signature seems to be really tricky!
												return "FORCE_REGENERATE", nil
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
												panic(fmt.Sprintf("Unexpected AST type: %T for %s.%s", v, pkg, iName))
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
	return "", errors.Newf(ctx, "interface %s not found in %s", iName, fullPath)
}

// resolvePackageDir asks `go list` where a package lives on disk so we can
// parse standard library and module packages that are not present in the local tree.
func resolvePackageDir(ctx context.Context, cwd, pkg string) (string, error) {
	cmd := exec.Command("go", "list", "-f", "{{.Dir}}", pkg)
	cmd.Dir = cwd
	output, err := cmd.Output()
	if err != nil {
		return "", errors.Wrapf(ctx, err, "resolve package directory for %s", pkg)
	}

	dir := strings.TrimSpace(string(output))
	if dir == "" {
		return "", errors.Newf(ctx, "empty package directory for %s", pkg)
	}
	if !isDir(dir) {
		return "", errors.Newf(ctx, "package directory does not exist: %s", dir)
	}
	return dir, nil
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
		return "interface{}"
	}

	return fmt.Sprintf("%v", field)
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
