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
	"time"
)

type MetaTime uint64

func MetaTimeNow() MetaTime {
	return MetaTimeSet(time.Now().UTC())
}

func MetaTimeSet(t time.Time) MetaTime {
	return MetaTime(uint64(t.Year())*1e13 + uint64(t.Month())*1e11 + uint64(t.Day())*1e9 +
		uint64(t.Hour()*1e7+t.Minute()*1e5+t.Second()*1e3) +
		uint64(t.Nanosecond()/1e6))
}

func (mt MetaTime) AddMillisecond(td int64) MetaTime {
	return MetaTimeSet(mt.Time().Add(time.Duration(td * 1e6)))
}

func (mt MetaTime) Add(ts string) MetaTime {
	td, _ := time.ParseDuration(ts)
	return MetaTimeSet(mt.Time().Add(td))
}

func (mt MetaTime) Format(fm string) string {

	if fm == "rfc3339" {
		fm = time.RFC3339
	}

	return mt.Time().Local().Format(fm)
}

func (mt MetaTime) Time() time.Time {

	mtu := uint64(mt)

	return time.Date(int(mtu/1e13), time.Month((mtu%1e13)/1e11), int((mtu%1e11)/1e9),
		int((mtu%1e9)/1e7), int((mtu%1e7)/1e5), int((mtu%1e5)/1e3),
		int(mtu%1e3)*1e6, time.UTC)
}
