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

package portutil

import (
	"fmt"
	"net"
	"testing"
)

func TestMain(t *testing.T) {

	start, end, err := FreeRange(10000, 200)
	if err != nil {
		t.Fatal(err.Error())
	}

	if start+200 != end {
		t.Fatal("Failed on FreeRange")
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", start))
	if err != nil {
		t.Fatal(err.Error())
	}
	defer ln.Close()

	if IsFree(start) == true {
		t.Fatal("Failed on IsFree")
	}

	_, err = Free(20000, 200)
	if err != nil {
		t.Fatal(err.Error())
	}

}
