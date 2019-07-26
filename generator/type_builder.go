package generator

import (
	"go/ast"
	"strconv"
	"strings"
)

type (
	TypeInfo struct {
		PkgPath string
		Name    string
		Fields  []Field
	}

	Field struct {
		Name string
		Type string
		Tags map[string]string
	}

	TypeInfoBuilder struct {
		PkgPath  string
		current  *TypeInfo
		AllTypes []TypeInfo
	}
)

func (builder *TypeInfoBuilder) StartNew(typeSpec *ast.TypeSpec) {

	builder.current = &TypeInfo{
		PkgPath: builder.PkgPath,
		Name:    typeSpec.Name.Name,
	}
}

func (builder *TypeInfoBuilder) AddFields(fieldList *ast.FieldList) {
	if builder.current == nil {
		return
	}
	fields := []Field{}
	for _, field := range fieldList.List {
		expr := field.Type
		if typ, ok := expr.(*ast.Ident); ok {
			field := Field{
				Name: field.Names[0].Name,
				Type: typ.Name,
				Tags: buildMap(field.Tag),
			}
			fields = append(fields, field)
		}
	}

	(*builder.current).Fields = fields
}

func (builder *TypeInfoBuilder) EndCurrent() {
	if builder.current != nil {
		builder.AllTypes = append(builder.AllTypes, *builder.current)
		builder.current = nil
	}
}

func buildMap(tag *ast.BasicLit) map[string]string {

	tags := map[string]string{}

	if tag == nil || tag.Value == "" {
		return tags
	}

	tagListRaw := tag.Value
	tagListRaw, _ = strconv.Unquote(tagListRaw)
	allTags := strings.Split(tagListRaw, " ")

	for _, tag := range allTags {
		kv := strings.Split(tag, ":")
		tags[kv[0]], _ = strconv.Unquote(kv[1])
	}

	return tags
}
