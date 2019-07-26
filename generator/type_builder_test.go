package generator

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildsATypeInfoObjectAfterAStartAndEndEvent(t *testing.T) {

	builder := TypeInfoBuilder{
		PkgPath: "pkg",
	}

	builder.StartNew(&ast.TypeSpec{Name: &ast.Ident{Name: "Entity"}})

	builder.EndCurrent()

	types := builder.AllTypes

	assert.Equal(t, len(types), 1)
	assert.Equal(t, types[0], TypeInfo{Name: "Entity", PkgPath: "pkg"})
}

func TestBuilderDoesNotHaveAnTypeInfoObjectWhenDoEndEventWasCalled(t *testing.T) {

	builder := TypeInfoBuilder{
		PkgPath: "pkg",
	}

	builder.StartNew(&ast.TypeSpec{Name: &ast.Ident{Name: "Entity"}})

	types := builder.AllTypes

	assert.Equal(t, len(types), 0)
}

func TestBuilderIgnoresAddFieldsEventThereWasNoStartEvent(t *testing.T) {

	builder := TypeInfoBuilder{
		PkgPath: "pkg",
	}

	builder.AddFields(&ast.FieldList{})

	types := builder.AllTypes

	assert.Equal(t, len(types), 0)
}

func TestBuilderAddFieldsToTypeInfoThatWasCreatedByAStartEvent(t *testing.T) {

	builder := TypeInfoBuilder{
		PkgPath: "pkg",
	}

	builder.StartNew(&ast.TypeSpec{Name: &ast.Ident{Name: "Entity"}})

	builder.AddFields(&ast.FieldList{
		List: []*ast.Field{
			&ast.Field{
				Names: []*ast.Ident{
					&ast.Ident{
						Name: "Field1",
					},
				},
				Tag:  &ast.BasicLit{},
				Type: &ast.Ident{Name: "string"},
			},
		},
	})

	builder.EndCurrent()
	types := builder.AllTypes

	assert.Equal(t, len(types), 1)
	assert.Equal(t, TypeInfo{Name: "Entity", PkgPath: "pkg", Fields: []Field{Field{Name: "Field1", Type: "string", Tags: map[string]string{}}}}, types[0])
}

func TestBuilderShouldParseTagToAKeyValueMap(t *testing.T) {

	builder := TypeInfoBuilder{
		PkgPath: "pkg",
	}

	builder.StartNew(&ast.TypeSpec{Name: &ast.Ident{Name: "Entity"}})

	builder.AddFields(&ast.FieldList{
		List: []*ast.Field{
			&ast.Field{
				Names: []*ast.Ident{
					&ast.Ident{
						Name: "Field1",
					},
				},
				Tag:  &ast.BasicLit{Value: "`key1:\"value1\" key2:\"value2\"`"},
				Type: &ast.Ident{Name: "string"},
			},
		},
	})

	builder.EndCurrent()
	types := builder.AllTypes

	assert.Equal(t, len(types), 1)
	assert.Equal(t, TypeInfo{Name: "Entity", PkgPath: "pkg", Fields: []Field{Field{Name: "Field1", Type: "string", Tags: map[string]string{"key1": "value1", "key2": "value2"}}}}, types[0])
}
