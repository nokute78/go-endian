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
	"fmt"
	"io"
	"reflect"
)

// write writes v to b.
func write(v reflect.Value, order ByteOrder, b []byte, index *int) error {
	var err error

	if !v.CanInterface() {
		// skip unexported field
		*index += sizeOfValue(v, false)
		return errCannotInterface
	}

	d := v.Interface()

	switch d.(type) {
	case uint8:
		b[*index] = d.(uint8)
		*index++
	case uint16:
		bs := make([]byte, 2)
		order.PutUint16(bs, d.(uint16))
		b[*index] = bs[0]
		b[*index+1] = bs[1]
		*index += 2

	case uint32:
		bs := make([]byte, 4)
		order.PutUint32(bs, d.(uint32))
		b[*index] = bs[0]
		b[*index+1] = bs[1]
		b[*index+2] = bs[2]
		b[*index+3] = bs[3]
		*index += 4

	case uint64:
		bs := make([]byte, 8)
		order.PutUint64(bs, d.(uint64))
		b[*index] = bs[0]
		b[*index+1] = bs[1]
		b[*index+2] = bs[2]
		b[*index+3] = bs[3]
		b[*index+4] = bs[4]
		b[*index+5] = bs[5]
		b[*index+6] = bs[6]
		b[*index+7] = bs[7]
		*index += 8
	default:
		switch v.Kind() {
		case reflect.Slice, reflect.Array:
			if v.Len() > 0 {
				if v.Index(0).Kind() == reflect.Uint8 {
					// byte slice / byte array
					length := v.Len()
					if order == BigEndian {
						for i := 0; i < length; i++ {
							b[*index+length-1-i] = byte(v.Index(i).Uint())
						}
					} else {
						for i := 0; i < length; i++ {
							b[*index+i] = byte(v.Index(i).Uint())
						}
					}
					*index += length
				} else {
					for i := 0; i < v.Len(); i++ {
						err := write(v.Index(i), order, b, index)
						if err != nil && err != errCannotInterface {
							return err
						}
					}
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
						*index += sizeOfValue(v.Field(i), true)
						continue
					} else if cnf.endian != Endian_Type_BLANK {
						var err error
						if cnf.endian == Endian_Type_BE {
							err = write(v.Field(i), BigEndian, b, index)
						} else {
							err = write(v.Field(i), LittleEndian, b, index)
						}
						if err != nil && err != errCannotInterface {
							return err
						}
						continue
					}
				}
				err := write(v.Field(i), order, b, index)
				if err != nil && err != errCannotInterface {
					return err
				}
			}
			return nil
		default:
			return fmt.Errorf("Not Supported %s", v.Kind())
		}
	}

	return err
}

// Write writes structured binary data from input into w.
func Write(w io.Writer, order ByteOrder, input interface{}) error {
	v := reflect.ValueOf(input)
	var vv reflect.Value

	switch v.Kind() {
	case reflect.Ptr:
		vv = reflect.Indirect(reflect.ValueOf(input))
	case reflect.Array, reflect.Slice, reflect.Struct:
		vv = reflect.ValueOf(input)
	default:
		return binary.Write(w, order, input)
	}

	barr := make([]byte, sizeOfValue(vv, true))
	index := 0
	err := write(vv, order, barr, &index)
	_, err = w.Write(barr)
	return err
}
