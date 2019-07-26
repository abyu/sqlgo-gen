package generator

import (
	"go/ast"

	"golang.org/x/tools/go/packages"
)

func LoadPackage(packageSourcePath string) []TypeInfo {
	cfg := &packages.Config{
		Mode: packages.LoadTypes | packages.NeedTypesInfo | packages.LoadSyntax,
	}

	pkgs, _ := packages.Load(cfg, packageSourcePath)
	builder := TypeInfoBuilder{PkgPath: pkgs[0].PkgPath}
	file := pkgs[0].Syntax

	for _, f := range file {
		ast.Inspect(f, dec(&builder))
	}

	return builder.AllTypes
}

func dec(builder *TypeInfoBuilder) func(ast.Node) bool {

	return func(node ast.Node) bool {

		if spec, ok := node.(*ast.TypeSpec); ok {
			builder.StartNew(spec)
		}
		if spec, ok := node.(*ast.FieldList); ok {
			builder.AddFields(spec)
			for i := 0; i < spec.NumFields(); i++ {
			}
			builder.EndCurrent()
		}
		return true
	}
}
