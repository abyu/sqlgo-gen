package generator

import (
	"reflect"
	"strconv"
	"strings"
)

func Parse(entity interface{}) *TypeInfo {

	entityType := reflect.ValueOf(entity).Type()
	fields := []Field{}

	for i := 0; i < entityType.NumField(); i++ {
		f := entityType.Field(i)
		field := Field{
			Name: f.Name,
			Type: f.Type.Name(),
			Tags: convertToMap(string(f.Tag)),
		}

		fields = append(fields, field)

	}

	return &TypeInfo{
		Name:    entityType.Name(),
		PkgPath: entityType.PkgPath(),
		Fields:  fields,
	}
}

func convertToMap(tag string) map[string]string {

	tags := map[string]string{}

	if tag == "" {
		return tags
	}

	tagListRaw := tag
	tagListRaw, _ = strconv.Unquote(tagListRaw)
	allTags := strings.Split(tagListRaw, " ")

	for _, tag := range allTags {
		kv := strings.Split(tag, ":")
		tags[kv[0]], _ = strconv.Unquote(kv[1])
	}

	return tags
}
