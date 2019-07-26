package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	SampleTestType struct {
		Field string
	}
)

var aTypeInUsedForTestsFile = "AType"
var aTypeInTestFile = "SampleTestType"

func TestShouldParseAPackageIntoAListOfTypeInfo(t *testing.T) {

	types := LoadPackage(".")
	typeNames := []string{}
	for _, typeInfo := range types {
		typeNames = append(typeNames, typeInfo.Name)

	}

	assert.NotContains(t, typeNames, aTypeInTestFile)
	assert.Contains(t, typeNames, aTypeInUsedForTestsFile)
}
