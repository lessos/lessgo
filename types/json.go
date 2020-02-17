// Copyright 2013-2018 lessgo Author, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	jsonTypelessFieldNameReg = regexp.MustCompile("^[a-zA-Z0-9_-]{1,50}$")
)

type JsonTypelessItem []*JsonTypelessField

func (it *JsonTypelessItem) Encode(indent string) ([]byte, error) {
	return json.MarshalIndent(it, "", indent)
}

func JsonTypelessItemDecode(src []byte) (JsonTypelessItem, error) {
	var item JsonTypelessItem
	err := json.Unmarshal(src, &item)
	return item, err
}

type JsonTypelessField struct {
	Name  string      `json:"name" toml:"name"`
	Value interface{} `json:"value,omitempty" toml:"value,omitempty"`
}

func (it *JsonTypelessField) val() ([]byte, string, int64, float64, bool) {

	bs, str, i64, f64, b := []byte{}, "", int64(0), float64(0), false

	switch it.Value.(type) {

	case bool:
		b = it.Value.(bool)

	case []byte:
		bs = it.Value.([]byte)

	case string:
		str = it.Value.(string)

	case int:
		i64 = int64(it.Value.(int))
	case int8:
		i64 = int64(it.Value.(int8))
	case int16:
		i64 = int64(it.Value.(int16))
	case int32:
		i64 = int64(it.Value.(int32))
	case int64:
		i64 = int64(it.Value.(int64))
	case uint:
		i64 = int64(it.Value.(uint))
	case uint8:
		i64 = int64(it.Value.(uint8))
	case uint16:
		i64 = int64(it.Value.(uint16))
	case uint32:
		i64 = int64(it.Value.(uint32))
	case uint64:
		i64 = int64(it.Value.(uint64))

	case float32:
		f64 = float64(it.Value.(float32))
	case float64:
		f64 = float64(it.Value.(float64))
	}

	return bs, str, i64, f64, b
}

func (it *JsonTypelessField) Bytes() []byte {
	bs, str, _, _, _ := it.val()
	if len(bs) > 0 {
		return bs
	}
	return []byte(str)
}

func (it *JsonTypelessField) String() string {
	bs, str, _, _, _ := it.val()
	if str != "" {
		return str
	}
	return string(bs)
}

func (it *JsonTypelessField) Bool() bool {
	_, _, _, _, b := it.val()
	return b
}

func (it *JsonTypelessField) Int64() int64 {
	_, str, i64, f64, _ := it.val()
	if i64 == 0 {
		if f64 != 0 {
			i64 = int64(f64)
		} else if str != "" {
			i64, _ = strconv.ParseInt(str, 10, 64)
		}
	}
	return i64
}

func (it *JsonTypelessField) Float64() float64 {
	_, str, i64, f64, _ := it.val()
	if f64 == 0 {
		if i64 != 0 {
			f64 = float64(i64)
		} else if str != "" {
			f64, _ = strconv.ParseFloat(str, 64)
		}
	}
	return f64
}

func (it *JsonTypelessField) IsArray() bool {

	switch it.Value.(type) {

	case []JsonTypelessItem:
		return true

	case []interface{}:
		return true
	}

	return false
}

func (it *JsonTypelessField) Array() []JsonTypelessItem {

	var items []JsonTypelessItem

	switch it.Value.(type) {
	case []JsonTypelessItem:
		return it.Value.([]JsonTypelessItem)

	case []interface{}:
		ls := it.Value.([]interface{})
		for _, v := range ls {
			if ls, ok := v.(map[string]interface{}); ok {
				var item JsonTypelessItem
				for k, v := range ls {
					item.Set(k, v)
				}
				items = append(items, item)
			}
		}

	}

	return items
}

func (it *JsonTypelessField) IsItem() bool {

	switch it.Value.(type) {

	case JsonTypelessItem:
		return true

	case map[string]interface{}:
		return true
	}

	return false
}

func (it *JsonTypelessField) Item() JsonTypelessItem {

	var item JsonTypelessItem

	switch it.Value.(type) {

	case JsonTypelessItem:
		return it.Value.(JsonTypelessItem)
	case map[string]interface{}:
		ls := it.Value.(map[string]interface{})
		for k, v := range ls {
			item.Set(k, v)
		}
	}

	return item
}

func (ls *JsonTypelessItem) Sort() {

	for i, v := range *ls {

		switch v.Value.(type) {

		case JsonTypelessItem:
			v2 := v.Value.(JsonTypelessItem)
			v2.Sort()
			(*ls)[i].Value = v2

		case []JsonTypelessItem:
			v2 := v.Value.([]JsonTypelessItem)
			for j, v22 := range v2 {
				v22.Sort()
				v2[j] = v22
			}
			(*ls)[i].Value = v2
		}
	}

	sort.Slice(*ls, func(i, j int) bool {
		return strings.Compare((*ls)[i].Name, (*ls)[j].Name) < 0
	})
}

func (ls *JsonTypelessItem) Get(name string) *JsonTypelessField {

	for _, v := range *ls {
		if v.Name == name {
			return v
		}
	}
	return nil
}

func (ls *JsonTypelessItem) Set(name string, value interface{}) {
	if prev := ls.Get(name); prev == nil {
		*ls = append(*ls, &JsonTypelessField{
			Name:  name,
			Value: value,
		})
	} else {
		prev.Value = value
	}
}

func (ls *JsonTypelessItem) ArraySet(name string, feed JsonTypelessItem) {

	for i, v := range *ls {

		if v.Name != name {
			continue
		}

		switch v.Value.(type) {
		case []interface{}:

			v2 := v.Value.([]JsonTypelessItem)
			v2 = append(v2, feed)

			(*ls)[i].Value = v2
		}

		return
	}

	ls.Set(name, []JsonTypelessItem{
		feed,
	})
}

func (ls *JsonTypelessItem) UnmarshalJSON(b []byte) error {

	var items map[string]interface{}
	if err := json.Unmarshal(b, &items); err != nil {
		return err
	}

	for k, v := range items {
		if !jsonTypelessFieldNameReg.MatchString(k) {
			return errors.New("invalid field name " + k)
		}

		*ls = append(*ls, &JsonTypelessField{
			Name:  k,
			Value: v,
		})
	}

	return nil
}

func (ls JsonTypelessItem) MarshalJSON() ([]byte, error) {

	var buf bytes.Buffer

	buf.WriteString("{")

	for i, v := range ls {

		if i != 0 {
			buf.WriteString(",")
		}

		key, err := json.Marshal(v.Name)
		if err != nil {
			return nil, err
		}

		val, err := json.Marshal(v.Value)
		if err != nil {
			return nil, err
		}

		buf.Write(key)
		buf.WriteString(":")
		buf.Write(val)
	}

	buf.WriteString("}")
	return buf.Bytes(), nil
}
