package environment

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	secretTagKey = "secret"
	secretMask   = "********"
)

func (env Environment) String() string {
	var b strings.Builder
	writeStruct(&b, reflect.ValueOf(env))
	return b.String()
}

func writeStruct(b *strings.Builder, value reflect.Value) {
	value = dereference(value)
	if !value.IsValid() {
		b.WriteString("<nil>")
		return
	}

	valueType := value.Type()
	b.WriteByte('{')
	for i := 0; i < value.NumField(); i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		field := valueType.Field(i)
		b.WriteString(field.Name)
		b.WriteByte(':')
		writeValue(b, field, value.Field(i))
	}
	b.WriteByte('}')
}

func writeValue(b *strings.Builder, field reflect.StructField, value reflect.Value) {
	value = dereference(value)
	if !value.IsValid() {
		b.WriteString("<nil>")
		return
	}

	if _, isSecret := field.Tag.Lookup(secretTagKey); isSecret && value.Kind() == reflect.String {
		if value.Len() > 0 {
			b.WriteString(secretMask)
		}
		return
	}

	if value.Kind() == reflect.Struct {
		writeStruct(b, value)
		return
	}

	fmt.Fprintf(b, "%+v", value.Interface())
}

func dereference(value reflect.Value) reflect.Value {
	for value.IsValid() && value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return reflect.Value{}
		}
		value = value.Elem()
	}
	return value
}
