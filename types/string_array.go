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
	"fmt"
	"regexp"
	"strings"
)

var (
	sarPatName        = regexp.MustCompile("^[a-z]{1}[a-z0-9-._]{4,50}$")
	sarErrNameLength  = errors.New("Length of the name must be between 5 and 50")
	sarErrNameInvalid = errors.New("Invalid Name")
)

type StringArray string

func (sa StringArray) String() string {
	return string(sa)
}

func (sa *StringArray) Set(namei interface{}) error {

	name := fmt.Sprintf("%v", namei)

	if len(name) < 4 || len(name) > 50 {
		return sarErrNameLength
	}

	if !labelPatName.MatchString(name) {
		return sarErrNameInvalid
	}

	if strings.Index(","+string(*sa)+",", ","+name+",") >= 0 {
		return nil
	}

	if len(string(*sa)) > 0 {
		*sa += StringArray("," + name)
	} else {
		*sa = StringArray(name)
	}

	return nil
}

func (sa *StringArray) Del(name interface{}) {

	names := fmt.Sprintf("%v", name)

	if i := strings.Index(","+string(*sa)+",", ","+names+","); i >= 0 {
		*sa = StringArray(strings.Trim(strings.Replace(","+string(*sa)+",", ","+names+",", ",", 1), ","))
	}
}

func (sa *StringArray) Has(names ...interface{}) bool {

	if len(names) < 1 {
		return false
	}

	sas := "," + string(*sa) + ","

	for _, v := range names {

		if strings.Index(sas, fmt.Sprintf(",%v,", v)) >= 0 {
			continue
		}

		return false
	}

	return true
}
