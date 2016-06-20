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

package json

import (
	"errors"
	"io/ioutil"
	"os"
)

func EncodeToFile(v interface{}, file, indent string) error {

	bs, err := Encode(v, indent)
	if err != nil {
		return err
	}

	if len(bs) < 2 {
		return errors.New("No Data Found")
	}

	fp, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()

	fp.Seek(0, 0)
	fp.Truncate(0)

	_, err = fp.Write(bs)

	return err
}

func DecodeFile(file string, v interface{}) error {

	fp, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fp.Close()

	bs, err := ioutil.ReadAll(fp)
	if err != nil {
		return err
	}

	return Decode(bs, v)
}
