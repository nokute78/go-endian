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
	"strings"
)

const (
	tagKeyName = "endian"
)

const (
	Endian_Type_BLANK = iota
	Endian_Type_LE
	Endian_Type_BE
)

// tagConfig represents StructTag.
//   "-"   : ignore the field
//   "skip": ignore but offset will be updated
//   "BE"  : the field is treated as big endian
//   "LE"  : the field is treated as little endian
type tagConfig struct {
	ignore bool
	skip   bool
	endian int
}

func parseStructTag(t reflect.StructTag) *tagConfig {
	s, ok := t.Lookup(tagKeyName)
	if !ok {
		return nil
	}
	ret := &tagConfig{}

	strs := strings.Split(s, ",")
	for _, v := range strs {
		switch v {
		case "-":
			ret.ignore = true
			return ret
		case "skip":
			ret.skip = true
			return ret
		case "BE":
			ret.endian = Endian_Type_BE
		case "LE":
			ret.endian = Endian_Type_LE
		}

	}
	return ret
}
