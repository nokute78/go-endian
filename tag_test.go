/*
   Copyright 2020 Takahiro Yamashita

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package endian

import (
	"reflect"
	"testing"
)

func TestParseStructTag(t *testing.T) {
	type A struct {
		Ignore bool      `endian:"-"`
		Skip   bool      `endian:"skip"`
		LE     ByteOrder `endian:"LE"`
		BE     ByteOrder `endian:"BE"`
	}

	a := A{}
	val := reflect.ValueOf(a)
	for i := 0; i < val.Type().NumField(); i++ {
		f := val.Type().Field(i)
		tagStr, ok := f.Tag.Lookup(tagKeyName)
		if !ok {
			t.Errorf("%d: Tag.Lookup is not ok", i)
			continue
		}

		cnf := parseStructTag(f.Tag)
		if cnf == nil {
			t.Errorf("%d: cnf is nil", i)
			continue
		}

		switch tagStr {
		case "-":
			if !cnf.ignore {
				t.Errorf("%d: key is - but ignore is false", i)
				continue
			}
		case "skip":
			if !cnf.skip {
				t.Errorf("%d: key is skip but skip is false", i)
				continue
			}
		case "LE":
			if cnf.endian != Endian_Type_LE {
				t.Errorf("%d: tag is LE but endian is not LittleEndian", i)
				continue
			}
		case "BE":
			if cnf.endian != Endian_Type_BE {
				t.Errorf("%d: tag is BE but endian is not BigEndian", i)
				continue
			}
		}
	}
}
