package generator

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/packages"
)

func LoadPackage() []TypeInfo {
	cfg := &packages.Config{
		Mode: packages.LoadTypes | packages.NeedTypesInfo | packages.LoadSyntax,
	}

	pkgs, err := packages.Load(cfg, "./models")
	builder := TypeInfoBuilder{PkgPath: pkgs[0].PkgPath}
	file := pkgs[0].Syntax

	ast.Inspect(file[0], dec(&builder))
	fmt.Println(err)

	return builder.AllTypes
	// fmt.Println(pkgs[0].Syntax)
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

type TypeInfo struct {
	PkgPath string
	Name    string
	Fields  []Field
}

type Field struct {
	Name string
	Type string
	Tags map[string]string
}

type TypeInfoBuilder struct {
	PkgPath  string
	current  *TypeInfo
	AllTypes []TypeInfo
}

func (builder *TypeInfoBuilder) StartNew(typeSpec *ast.TypeSpec) {

	builder.current = &TypeInfo{
		PkgPath: builder.PkgPath,
		Name:    typeSpec.Name.Name,
	}
}

func (builder *TypeInfoBuilder) AddFields(fieldList *ast.FieldList) {

	fields := []Field{}
	for _, field := range fieldList.List {
		expr := field.Type
		if typ, ok := expr.(*ast.Ident); ok {
			field := Field{
				Name: field.Names[0].Name,
				Type: typ.Name,
				Tags: buildMap(field.Tag.Value),
			}
			fields = append(fields, field)
		}
	}
	if builder.current != nil {
		(*builder.current).Fields = fields
	}
}

func (builder *TypeInfoBuilder) EndCurrent() {
	if builder.current != nil {
		builder.AllTypes = append(builder.AllTypes, *builder.current)
		builder.current = nil
	}
}

func buildMap(tagListRaw string) map[string]string {
	tags := map[string]string{}

	tagListRaw = strings.ReplaceAll(tagListRaw, "`", "")
	tagListRaw = strings.ReplaceAll(tagListRaw, `"`, "")
	allTags := strings.Split(tagListRaw, " ")

	for _, tag := range allTags {
		kv := strings.Split(tag, ":")
		tags[kv[0]] = kv[1]
	}

	return tags
}
