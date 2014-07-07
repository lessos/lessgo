package base

import (
	"fmt"
	"reflect"
	"time"
)

type Field struct {
	value     reflect.Value
	valueType reflect.Type
}

type Entry struct {
	Fields map[string]*Field
}

func (e *Entry) Field(fieldName string) *Field {

	field, ok := e.Fields[fieldName]
	if ok {
		return field
	}

	return &Field{}
}

func (f *Field) String() string {

	if f.value.Interface() == nil {
		return ""
	}

	vv := reflect.ValueOf(f.value.Interface())

	switch f.valueType.Kind() {
	case reflect.Slice:
		if f.valueType.Elem().Kind() == reflect.Uint8 {
			return string(f.value.Interface().([]byte))
		}
	case reflect.String:
		return vv.String()
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", vv.Int())
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", vv.Uint())
	}

	return ""
}

func (f *Field) Int() int64 {

	if f.value.Interface() == nil {
		return 0
	}

	vv := reflect.ValueOf(f.value.Interface())

	switch f.valueType.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return vv.Int()
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(vv.Uint())
	}

	return 0
}

func (f *Field) Uint() uint64 {

	if f.value.Interface() == nil {
		return 0
	}

	vv := reflect.ValueOf(f.value.Interface())

	switch f.valueType.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint64(vv.Int())
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return vv.Uint()
	}

	return 0
}

func (f *Field) Float() float64 {

	if f.value.Interface() == nil {
		return 0
	}

	vv := reflect.ValueOf(f.value.Interface())

	switch f.valueType.Kind() {
	case reflect.Float32, reflect.Float64:
		return vv.Float()
	}

	return 0
}

func (f *Field) TimeParse(format string) time.Time {

	if f.value.Interface() == nil {
		return time.Now().In(TimeZone)
	}

	vv := reflect.ValueOf(f.value.Interface())

	timeString := ""
	switch f.valueType.Kind() {
	case reflect.Slice:
		if f.valueType.Elem().Kind() == reflect.Uint8 {
			timeString = string(f.value.Interface().([]byte))
		}
	case reflect.String:
		timeString = vv.String()
	}

	return TimeParse(timeString, format)
}
