package main

import (
	"github.com/abyu/sqlxx/generator"
	"github.com/abyu/sqlxx/repository"
)

func main() {
	// anId := int64(10)
	// theEntity := getById(anId)

	// repository.GenerateStruct("EntityRepository", models.AnEntity{})

	typeInfo := generator.LoadPackage()

	for _, info := range typeInfo {
		repository.GenerateStruct(info.Name, info)
	}
	// fmt.Println(theEntity)
}
