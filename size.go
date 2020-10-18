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
)

func sizeOfValueRecursive(c *int, v reflect.Value, structtag bool) {
	switch v.Kind() {
	case reflect.Struct:
		if structtag {
			for i := 0; i < v.Type().NumField(); i++ {
				f := v.Type().Field(i)
				cnf := parseStructTag(f.Tag)
				if cnf != nil && cnf.ignore {
					continue
				}
				sizeOfValueRecursive(c, v.Field(i), structtag)
			}
		} else {
			for i := 0; i < v.NumField(); i++ {
				sizeOfValueRecursive(c, v.Field(i), structtag)
			}
		}
	case reflect.Array, reflect.Slice:
		if v.Len() == 0 {
			return
		}
		var elemSize int
		sizeOfValueRecursive(&elemSize, v.Index(0), structtag)
		*c += (elemSize * v.Len())
	default:
		/* other types */
		*c += v.Type().Bits() / 8
	}
}

func sizeOfValue(v reflect.Value, structtag bool) (ret int) {
	sizeOfValueRecursive(&ret, v, structtag)
	return ret
}
