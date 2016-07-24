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
	"testing"
)

func TestNameIdentifier(t *testing.T) {

	as_valid := []NameIdentifier{
		"internal/service/webserver",
		"example.com/service/webserver",
		"sub-domain.example.com/service/webserver",
	}

	for _, v := range as_valid {

		if err := v.Valid(); err != nil {
			t.Fatal("Fatal on Valid: " + err.Error())
		}
	}

	as_invalid := []NameIdentifier{
		"internal",
		"!",
		"",
		"@abc#",
	}

	for _, v := range as_invalid {

		if err := v.Valid(); err == nil {
			t.Fatal("Fatal on Valid")
		}
	}
}
