// Copyright 2015 lessOS.com, All rights reserved.
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

package httpsrv

import (
	"html/template"
	"reflect"
	"strings"
	"time"
)

var templateFuncs = map[string]interface{}{
	"eq": tfEqual,
	// Skips sanitation on the parameter.  Do not use with dynamic data.
	"raw": func(text string) template.HTML {
		return template.HTML(text)
	},
	// Returns a copy of the string s with the old replaced by new
	"replace": func(s, old, new string) string {
		return strings.Replace(s, old, new, -1)
	},
	// Format a date according to the application's default date(time) format.
	"date": func(t time.Time) string {
		return t.Format("2006-01-02")
	},
	"datetime": func(t time.Time) string {
		//t, _ := time.Parse("2006-01-02 15:04:05.000 -0700", fmttime)
		return t.Format("2006-01-02 15:04")
	},
	// upper returns a copy of the string s with all Unicode letters mapped to their upper case
	"upper": func(s string) string {
		return strings.ToUpper(s)
	},
	// lower returns a copy of the string s with all Unicode letters mapped to their lower case
	"lower": func(s string) string {
		return strings.ToLower(s)
	},
	// Perform a message look-up for the given locale and message using the given arguments
	"T": func(lang map[string]interface{}, msg string, args ...interface{}) string {
		return i18nTranslate(lang["LANG"].(string), msg, args...)
	},
	// "set": func(renderArgs map[string]interface{}, key string, value interface{}) {
	// 	renderArgs[key] = value
	// },
}

// Equal is a helper for comparing value equality, following these rules:
//  - Values with equivalent types are compared with reflect.DeepEqual
//  - int, uint, and float values are compared without regard to the type width.
//    for example, Equal(int32(5), int64(5)) == true
//  - strings and byte slices are converted to strings before comparison.
//  - else, return false.
func tfEqual(a, b interface{}) bool {
	if reflect.TypeOf(a) == reflect.TypeOf(b) {
		return reflect.DeepEqual(a, b)
	}
	switch a.(type) {
	case int, int8, int16, int32, int64:
		switch b.(type) {
		case int, int8, int16, int32, int64:
			return reflect.ValueOf(a).Int() == reflect.ValueOf(b).Int()
		}
	case uint, uint8, uint16, uint32, uint64:
		switch b.(type) {
		case uint, uint8, uint16, uint32, uint64:
			return reflect.ValueOf(a).Uint() == reflect.ValueOf(b).Uint()
		}
	case float32, float64:
		switch b.(type) {
		case float32, float64:
			return reflect.ValueOf(a).Float() == reflect.ValueOf(b).Float()
		}
	case string:
		switch b.(type) {
		case []byte:
			return a.(string) == string(b.([]byte))
		}
	case []byte:
		switch b.(type) {
		case string:
			return b.(string) == string(a.([]byte))
		}
	}
	return false
}

// NotEqual evaluates the comparison a != b
func tfNotEqual(a, b interface{}) bool {
	return !tfEqual(a, b)
}

// LessThan evaluates the comparison a < b
func tfLessThan(a, b interface{}) bool {

	switch a.(type) {
	case int, int8, int16, int32, int64:
		switch b.(type) {
		case int, int8, int16, int32, int64:
			return reflect.ValueOf(a).Int() < reflect.ValueOf(b).Int()
		}
	case uint, uint8, uint16, uint32, uint64:
		switch b.(type) {
		case uint, uint8, uint16, uint32, uint64:
			return reflect.ValueOf(a).Uint() < reflect.ValueOf(b).Uint()
		}
	case float32, float64:
		switch b.(type) {
		case float32, float64:
			return reflect.ValueOf(a).Float() < reflect.ValueOf(b).Float()
		}
	}

	return false
}

// LessEqual evaluates the comparison a <= b
func tfLessEqual(a, b interface{}) bool {
	return tfLessThan(a, b) || tfEqual(a, b)
}

// GreaterThan evaluates the comparison a > b
func tfGreaterThan(a, b interface{}) bool {
	return !tfLessEqual(a, b)
}

// GreaterEqual evaluates the comparison a >= b
func tfGreaterEqual(a, b interface{}) bool {
	return !tfLessThan(a, b)
}
