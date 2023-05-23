package struct_to_map

import (
	"reflect"
	"strings"
)

func GetTagByField(field reflect.StructField, fieldTag string) string {
	tag := field.Tag.Get(fieldTag)
	if tag == "" || tag == "-" {
		return field.Name
	}
	tag = strings.Split(tag, ",")[0]
	if tag == "" {
		return ""
	}
	return tag
}
