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
	"strings"
	"testing"
)

func TestKvPair(t *testing.T) {

	ls := KvPairs{}

	if err := ls.Set("", "value"); err == nil {
		t.Fatal("Failed on Set")
	}

	key101 := strings.Repeat("a", 101)
	if err := ls.Set(key101, "value"); err == nil {
		t.Fatal("Failed on Set")
	}

	ls.Set("aaaa", "value aaaa")
	ls.Set("bbbb", "value bbbb")
	ls.Set("cccc", "value cccc")
	ls.Set("dddd", "value dddd")
	ls.Set("eeee", "value eeee")

	if _, hit := ls.Get("aaaa"); !hit {
		t.Fatal("Failed on Get")
	}

	if _, hit := ls.Get("0000"); hit {
		t.Fatal("Failed on Get")
	}

	ls.Del("aaaa")
	if _, hit := ls.Get("aaaa"); hit {
		t.Fatal("Failed on Get")
	}

	ls.Del("cccc")
	if _, hit := ls.Get("cccc"); hit {
		t.Fatal("Failed on Get")
	}

	ls.Del("eeee")
	if _, hit := ls.Get("eeee"); hit {
		t.Fatal("Failed on Get")
	}

	ls2 := KvPairs{}

	ls2.Set("bbbb", "value bbbb")
	ls2.Set("dddd", "value dddd")

	if ls.Equal(ls2) == false {
		t.Fatal("Failed on Equal")
	}

	ls2.Set("dddd", "value changed")
	if ls.Equal(ls2) == true {
		t.Fatal("Failed on Equal")
	}
}
