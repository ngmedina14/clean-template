package common

import (
	"reflect"
	"strings"
)

func IsPointerToStruct(pointer interface{}) bool {
	return reflect.TypeOf(pointer).Kind() == reflect.Ptr && reflect.Indirect(reflect.ValueOf(pointer)).Kind() == reflect.Struct
}

func IsEmptyOrWhitespace(s string) bool {
	return strings.TrimSpace(s) == ""
}
