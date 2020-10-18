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

package endian_test

import (
	"bytes"
	"encoding/binary"
	"github.com/nokute78/go-endian"
	"io"
	"testing"
)

func TestReadPrimitive(t *testing.T) {
	// uint8
	br := bytes.NewReader([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08})
	var u8 uint8
	if err := endian.Read(br, endian.LittleEndian, &u8); err != nil {
		t.Errorf("error %s\n", err)
	}
	_, err := br.Seek(0, io.SeekStart)
	if err != nil {
		t.Errorf("%s", err)
	}
	if u8 != 0 {
		t.Errorf("uint8:given %x expect 0", u8)
	}

	// uint16
	var u16 uint16
	if err := endian.Read(br, endian.LittleEndian, &u16); err != nil {
		t.Errorf("error %s\n", err)
	}
	_, err = br.Seek(0, io.SeekStart)
	if err != nil {
		t.Errorf("%s", err)
	}
	if u16 != 0x100 {
		t.Errorf("uint16:given %x expect 0", u16)
	}

	// uint32
	var u32 uint32
	if err := endian.Read(br, endian.LittleEndian, &u32); err != nil {
		t.Errorf("error %s\n", err)
	}
	_, err = br.Seek(0, io.SeekStart)
	if err != nil {
		t.Errorf("%s", err)
	}
	if u32 != 0x03020100 {
		t.Errorf("uint32:given %x expect 0", u32)
	}

	// uint64
	var u64 uint64
	if err := endian.Read(br, endian.LittleEndian, &u64); err != nil {
		t.Errorf("error %s\n", err)
	}
	_, err = br.Seek(0, io.SeekStart)
	if err != nil {
		t.Errorf("%s", err)
	}
	if u64 != 0x0706050403020100 {
		t.Errorf("uint64:given %x expect 0", u64)
	}
}

func TestReadStruct(t *testing.T) {
	type Sample struct {
		Header   byte
		Reserved uint16
		Rev      byte
		Data     [4]byte
	}
	s := Sample{}
	br := bytes.NewReader([]byte{0x7f, 0xff, 0xff, 0x51, 0xaa, 0xbb, 0xcc, 0xdd})
	if err := endian.Read(br, endian.LittleEndian, &s); err != nil {
		t.Fatalf("error:%s", err)
	}

	if s.Header != 0x7f {
		t.Errorf("header: given=%x expect=%x", s.Header, 0x7f)
	}
	if s.Reserved != 0xffff {
		t.Errorf("reserved: given=%x expect=%x", s.Reserved, 0xffff)
	}
	if s.Rev != 0x51 {
		t.Errorf("rev: given=%x expect=%x", s.Rev, 0x51)
	}

	expect := []byte{0xaa, 0xbb, 0xcc, 0xdd}
	if len(s.Data) != len(expect) {
		t.Errorf("Data size error")
	}
	if bytes.Compare(expect, s.Data[:]) != 0 {
		t.Errorf("Data: given=%v expect=%v", s.Data, expect)
	}
}

func TestStructTag(t *testing.T) {
	type Sample struct {
		Reserved byte `endian:"-"` // ignored
		Val      byte
	}

	s := Sample{}
	br := bytes.NewReader([]byte{0xff, 0xaa})
	if err := endian.Read(br, endian.LittleEndian, &s); err != nil {
		t.Fatalf("error:%s\n", err)
	}

	// ignored field
	if s.Reserved != 0x00 {
		t.Errorf("given=0x%x expect=0x00", s.Reserved)
	}

	if s.Val != 0xff {
		t.Errorf("given=0x%x expect=0xaa", s.Val)
	}

	// skip case
	type Sample2 struct {
		Reserved byte `endian:"skip"` // skip
		Val      byte
	}

	s2 := Sample2{}
	br = bytes.NewReader([]byte{0xff, 0xaa})
	if err := endian.Read(br, endian.LittleEndian, &s2); err != nil {
		t.Fatalf("error:%s\n", err)
	}

	// ignored field
	if s2.Reserved != 0x00 {
		t.Errorf("given=0x%x expect=0x0", s2.Reserved)
	}
	if s2.Val != 0xaa {
		t.Errorf("given=0x%x expect=0xaa", s2.Val)
	}
}

func TestReadArray(t *testing.T) {
	input := bytes.NewBuffer([]byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff})
	var b [6]byte

	if err := endian.Read(input, endian.BigEndian, &b); err != nil {
		t.Fatalf("bit.Read:%s", err)
	}

	expect := []byte{0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa}
	if bytes.Compare(b[:], expect) != 0 {
		t.Fatalf("mismatch: given=%x expect=%x", b, expect)
	}
}

func TestMixedEndian(t *testing.T) {
	type Data struct {
		F1 uint32
		F2 uint16
		F3 uint16
		F4 uint16  `endian:"BE"`
		F5 [6]byte `endian:"BE"`
	}

	d := Data{}
	b := bytes.NewBuffer([]byte{0x57, 0xab, 0xdf, 0x5d, 0xa1, 0xdf, 0xaa, 0x4e, 0x96, 0xb5, 0x3a, 0x5f, 0xe7, 0x66, 0x92, 0x65})
	if err := endian.Read(b, endian.LittleEndian, &d); err != nil {
		t.Fatalf("endian.Read:%s", err)
	}

	if d.F1 != 0x5ddfab57 {
		t.Errorf("F1 mismatch:given=0x%x expect=0x%x", d.F1, 0x5ddfab57)
	}
	if d.F2 != 0xdfa1 {
		t.Errorf("F2 mismatch:given=0x%x expect=0x%x", d.F2, 0xdfa1)
	}
	if d.F3 != 0x4eaa {
		t.Errorf("F3 mismatch:given=0x%x expect=0x%x", d.F3, 0x4eaa)
	}

	if d.F4 != 0x96b5 {
		t.Errorf("F4 mismatch:given=0x%x expect=0x%x", d.F4, 0x96b5)
	}

	expect := []byte{0x65, 0x92, 0x66, 0xe7, 0x5f, 0x3a}
	if bytes.Compare(d.F5[:], expect) != 0 {
		t.Errorf("F5 mismatch:given=%x expect=%x", d.F5, expect)
	}
}

func BenchmarkReadStruct(b *testing.B) {
	type Sample struct {
		Header   byte
		Reserved [2]byte
		Id       byte
		Data     [4]byte
	}
	s := Sample{}
	br := bytes.NewReader([]byte{0x7f, 0xff, 0xff, 0x51, 0xaa, 0xbb, 0xcc, 0xdd})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := br.Seek(0, io.SeekStart)
		if err != nil {
			b.Fatalf("error:%s", err)
		}
		err = endian.Read(br, endian.LittleEndian, &s)
		if err != nil {
			b.Fatalf("error:%s", err)
		}
	}
}

func BenchmarkBinaryReadStruct(b *testing.B) {
	type Sample struct {
		Header   byte
		Reserved [2]byte
		Id       byte
		Data     [4]byte
	}
	s := Sample{}
	br := bytes.NewReader([]byte{0x7f, 0xff, 0xff, 0x51, 0xaa, 0xbb, 0xcc, 0xdd})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := br.Seek(0, io.SeekStart)
		if err != nil {
			b.Fatalf("error:%s", err)
		}
		err = binary.Read(br, endian.LittleEndian, &s)
		if err != nil {
			b.Fatalf("error:%s", err)
		}
	}
}

func TestUnexportedField(t *testing.T) {
	type S struct {
		A byte
		_ byte
		B byte
	}

	var s S
	buf := bytes.NewBuffer([]byte{0xaa, 0xfb, 0xcc})
	if err := endian.Read(buf, endian.LittleEndian, &s); err != nil {
		t.Errorf("err=%s", err)
	}
	var expect S
	buf = bytes.NewBuffer([]byte{0xaa, 0xfb, 0xcc})
	if err := binary.Read(buf, binary.LittleEndian, &expect); err != nil {
		t.Errorf("err=%s", err)
	}

	if s.A != expect.A {
		t.Errorf("s.A mistmach: given=%x expect=%x", s.A, expect.A)
	}
	if s.B != expect.B {
		t.Errorf("s.B mistmach: given=%x expect=%x", s.B, expect.B)
	}
}
