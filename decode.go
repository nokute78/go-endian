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
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
)

var errCannotInterface = errors.New("CanInterface returns false")

// read reads from b and fill v.
func read(b []byte, order ByteOrder, v reflect.Value, index *int) error {
	var val reflect.Value
	if !v.CanInterface() {
		// skip unexported field
		*index += sizeOfValue(v, false)
		return errCannotInterface
	}
	d := v.Interface()

	switch d.(type) {
	case uint8:
		val = reflect.ValueOf(b[*index])
		*index += 1
	case uint16:
		val = reflect.ValueOf(order.Uint16(b[*index : *index+2]))
		*index += 2
	case uint32:
		val = reflect.ValueOf(order.Uint32(b[*index : *index+4]))
		*index += 4
	case uint64:
		val = reflect.ValueOf(order.Uint64(b[*index : *index+8]))
		*index += 8
	default: /* other data types */
		switch v.Kind() {
		case reflect.Array, reflect.Slice:
			if v.Len() > 0 {
				if v.Index(0).Kind() == reflect.Uint8 {
					length := v.Len()
					ret := b[*index : *index+length]
					if order == BigEndian {
						for i := 0; i < length; i++ {
							if v.Index(i).CanSet() {
								// workaround! binary.Read doesn't support []byte in BigEndian
								v.Index(i).Set(reflect.ValueOf(ret[length-1-i]))
							}
						}
					} else {
						for i := 0; i < length; i++ {
							if v.Index(i).CanSet() {
								v.Index(i).Set(reflect.ValueOf(ret[i]))
							}
						}
					}
					*index += length
					return nil
				} else {
					for i := 0; i < v.Len(); i++ {
						err := read(b, order, v.Index(i), index)
						if err != nil && err != errCannotInterface {
							return err
						}
					}
					return nil
				}
			}
		case reflect.Struct:
			for i := 0; i < v.Type().NumField(); i++ {
				f := v.Type().Field(i)
				cnf := parseStructTag(f.Tag)
				if cnf != nil {
					/* struct tag is defined */
					if cnf.ignore {
						continue
					} else if cnf.skip {
						/* only updates offset. not fill. */
						*index = sizeOfValue(v.Field(i), true)
						continue
					} else if cnf.endian != nil {
						err := read(b, cnf.endian, v.Field(i), index)
						if err != nil && err != errCannotInterface {
							return err
						}
						continue
					}
				}
				err := read(b, order, v.Field(i), index)
				if err != nil && err != errCannotInterface {
					return err
				}
			}
			return nil
		default:
			return fmt.Errorf("Not Supported %s", v.Kind())
		}
	}

	// primitives
	if v.CanSet() {
		v.Set(val)
	} else {
		return fmt.Errorf("can not set %v\n", v)
	}
	return nil
}

// Read reads structured binary data from i into data.
// Data must be a pointer to a fixed-size value.
// Not exported struct field is ignored.
//   Supports StructTag.
//       `bit:"skip"` : ignore the field. Skip X bits which is the size of the field. It is useful for reserved field.
//       `bit:"-"`    : ignore the field. Offset is not changed.
func Read(r io.Reader, order ByteOrder, data interface{}) error {
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Ptr:
		c := sizeOfValue(reflect.Indirect(v), true)
		barr := make([]byte, c)
		n, err := r.Read(barr)
		if err != nil {
			return err
		} else if n != c {
			return fmt.Errorf("endian.Read:short read, expect=%d byte, read=%d byte", c, n)
		}
		index := 0
		err = read(barr, order, reflect.Indirect(v), &index)
		if err != io.EOF && err != errCannotInterface {
			return err
		}
	default:
		return binary.Read(r, order, data)
	}
	return nil
}
