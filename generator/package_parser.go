package generator

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/packages"
)

func LoadPackage(packageSourcePath string) []TypeInfo {
	cfg := &packages.Config{
		Mode: packages.LoadTypes | packages.NeedTypesInfo | packages.LoadSyntax,
	}

	pkgs, err := packages.Load(cfg, packageSourcePath)
	builder := TypeInfoBuilder{PkgPath: pkgs[0].PkgPath}
	file := pkgs[0].Syntax

	ast.Inspect(file[0], dec(&builder))
	fmt.Println(err)

	return builder.AllTypes
}

func dec(builder *TypeInfoBuilder) func(ast.Node) bool {

	return func(node ast.Node) bool {
		// fmt.Printf("--- %+v => %+v ---\n", reflect.TypeOf(node), node)

		if spec, ok := node.(*ast.TypeSpec); ok {
			builder.StartNew(spec)
			// fmt.Printf("[%v]", spec.Name)
		}
		// fmt.Print("{")
		if spec, ok := node.(*ast.FieldList); ok {
			builder.AddFields(spec)
			for i := 0; i < spec.NumFields(); i++ {
				// fmt.Printf("%v, ", spec.List[i].Tag.Value)
			}
			builder.EndCurrent()
		}
		// fmt.Print("}")
		return true
	}
}
