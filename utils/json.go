package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
)

const (
	jsonMergeDepth = 32
)

func JsonDecode(src, rs interface{}) (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = errors.New("json: invalid format")
		}
	}()

	var bf []byte
	switch src.(type) {
	case string:
		bf = []byte(src.(string))
	case []byte:
		bf = src.([]byte)
	default:
		panic("invalid format")
	}

	if err = json.Unmarshal(bf, &rs); err != nil {
		return err
	}

	return nil
}

func JsonEncode(rs interface{}) (str string, err error) {

	rb, err := json.Marshal(rs)

	if err == nil {
		str = string(rb)
	}

	return
}

func JsonEncodeBytes(rs interface{}) ([]byte, error) {
	return json.Marshal(rs)
}

func JsonEncodeIndent(rs interface{}, indent string) ([]byte, error) {
	return json.MarshalIndent(rs, "", indent)
}

func JsonIndent(src, indent string) ([]byte, error) {

	var out bytes.Buffer
	err := json.Indent(&out, []byte(src), "", indent)

	return out.Bytes(), err
}

// Merge recursively merges the src and dst maps. Key conflicts are resolved by
// preferring src, or recursively descending, if both src and dst are maps.
//
// Refer
// 	https://github.com/peterbourgon/mergemap
func JsonMerge(dst, src map[string]interface{}) map[string]interface{} {
	return jsMerge(dst, src, 0)
}

func jsMerge(dst, src map[string]interface{}, depth int) map[string]interface{} {

	if depth > jsonMergeDepth {
		return dst
		// panic("too deep!")
	}

	for key, srcVal := range src {

		if dstVal, ok := dst[key]; ok {

			srcMap, srcMapOk := jsMapify(srcVal)
			dstMap, dstMapOk := jsMapify(dstVal)

			if srcMapOk && dstMapOk {
				srcVal = jsMerge(dstMap, srcMap, depth+1)
			}
		}

		dst[key] = srcVal
	}

	return dst
}

func jsMapify(i interface{}) (map[string]interface{}, bool) {

	value := reflect.ValueOf(i)

	if value.Kind() == reflect.Map {

		m := map[string]interface{}{}

		for _, k := range value.MapKeys() {
			m[k.String()] = value.MapIndex(k).Interface()
		}

		return m, true
	}

	return map[string]interface{}{}, false
}
