package base

import (
	"../../../utils"
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

// Json returns the map that marshals from the reply bytes as json in response .
func (f *Field) Json(v interface{}) error {
	return utils.JsonDecode(f.String(), v)
}

func (f *Field) Int8() int8 {
	return int8(f.Int64())
}

func (f *Field) Int16() int16 {
	return int16(f.Int64())
}

func (f *Field) Int32() int32 {
	return int32(f.Int64())
}

func (f *Field) Int64() int64 {

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

func (f *Field) Int() int64 {
	return f.Int64()
}

func (f *Field) Uint8() uint8 {
	return uint8(f.Uint64())
}

func (f *Field) Uint16() uint16 {
	return uint16(f.Uint64())
}

func (f *Field) Uint32() uint32 {
	return uint32(f.Uint64())
}

func (f *Field) Uint64() uint64 {

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

func (f *Field) Uint() uint64 {
	return f.Uint64()
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
