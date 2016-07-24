// Copyright 2013-2016 lessgo Author, All rights reserved.
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
	"errors"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lessos/lessgo/crypto/idhash"
)

// NameIdentifier is restricted to lowercase URI unreserved characters
// defined in RFC3986 (http://tools.ietf.org/html/rfc3986#section-2.3)
//
// Format:
//  cannot be an empty string
//  must begin and end with an alphanumeric character.
//  include two or more components that are delimited by "/" characters
//  will match the RE2 regular expression: ^[a-z0-9]+([-._~/][a-z0-9]+)*$
//
// Examples:
//  internal/service/webserver
//  example.com/service/webserver
//  sub-domain.example.com/service/webserver
type NameIdentifier string

var (
	ni_re2        = regexp.MustCompile("^[a-z0-9]+([-._~/][a-z0-9]+)*$")
	ni_hashed_re2 = regexp.MustCompile("[0-9a-f]{4,32}$")
)

func NewNameIdentifier(str string) NameIdentifier {
	return NameIdentifier(strings.Trim(filepath.Clean(str), "/"))
}

func (ni NameIdentifier) String() string {
	return string(ni)
}

func (ni NameIdentifier) HashToString(length int) string {

	if length > 32 {
		length = 32
	} else if length < 4 {
		length = 4
	}

	return idhash.HashToHexString([]byte(ni), length)
}

func (ni NameIdentifier) IsHashed() bool {

	if ni_hashed_re2.MatchString(string(ni)) {
		return true
	}

	return false
}

func (ni NameIdentifier) Valid() error {

	nis := strings.Trim(filepath.Clean(string(ni)), "/")

	if len(nis) < 1 {
		return errors.New("NameIdentifier cannot be an empty string")
	}

	if strings.IndexByte(nis, '/') < 0 {
		return errors.New("NameIdentifier at least include one '/' character")
	}

	if !ni_re2.MatchString(nis) {
		return errors.New("Invalid NameIdentifier")
	}

	return nil
}
