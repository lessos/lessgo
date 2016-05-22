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
	"math"
	"testing"
	"time"
)

func TestTime(t *testing.T) {

	mt := MetaTimeNow()
	mn := time.Now().UTC()

	if math.Abs(float64(mt.Time().Unix()-mn.Unix())) > 1 {
		t.Fatal("Failed on MetaTimeNow")
	}

	if (mt.AddMillisecond(3600000).Time().Unix() - mt.Time().Unix()) != 3600 {
		t.Fatal("Failed on AddMillisecond")
	}

	if (mt.Add("+3600s").Time().Unix() - mt.Time().Unix()) != 3600 {
		t.Fatal("Failed on Add")
	}
}

func Benchmark_Time_MetaTimeNow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MetaTimeNow()
	}
}

func Benchmark_Time_Time(b *testing.B) {

	mt := MetaTimeNow()
	for i := 0; i < b.N; i++ {
		mt.Time()
	}
}

func Benchmark_Time_Format(b *testing.B) {

	mt := MetaTimeNow()
	for i := 0; i < b.N; i++ {
		mt.Format("rfc3339")
	}
}

func Benchmark_Time_AddMillisecond(b *testing.B) {

	mt := MetaTimeNow()

	for i := 0; i < b.N; i++ {
		mt.AddMillisecond(3600000)
	}
}

func Benchmark_Time_Add(b *testing.B) {

	mt := MetaTimeNow()

	for i := 0; i < b.N; i++ {
		mt.Add("+3600s")
	}
}
