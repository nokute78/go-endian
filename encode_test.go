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
	"testing"
)

func TestWritePrimitive(t *testing.T) {
	// uint8
	buf := bytes.NewBuffer([]byte{})
	var u8 uint8 = 0xbb
	if err := endian.Write(buf, endian.LittleEndian, u8); err != nil {
		t.Errorf("endian.Write error %s", err)
	}
	ret, err := buf.ReadByte()
	if err != nil {
		t.Errorf("ReadByte err=%s", err)
	} else if ret != u8 {
		t.Errorf("mismatch given 0x%x expect 0x%x", ret, u8)
	}

	// uint16
	buf.Reset()
	var u16 uint16 = 0xbbee
	if err := endian.Write(buf, endian.LittleEndian, u16); err != nil {
		t.Errorf("endian.Write error %s", err)
	}
	ret16 := endian.LittleEndian.Uint16(buf.Bytes())
	if ret16 != u16 {
		t.Errorf("mismatch given 0x%x expect 0x%x", ret16, u16)
	}

	// uint32
	buf.Reset()
	var u32 uint32 = 0xbbeeccff
	if err := endian.Write(buf, endian.LittleEndian, u32); err != nil {
		t.Errorf("endian.Write error %s", err)
	}
	ret32 := endian.LittleEndian.Uint32(buf.Bytes())
	if ret32 != u32 {
		t.Errorf("mismatch given 0x%x expect 0x%x", ret32, u32)
	}

	// uint64
	buf.Reset()
	var u64 uint64 = 0xbbeeccff00112233
	if err := endian.Write(buf, endian.LittleEndian, u64); err != nil {
		t.Errorf("endian.Write error %s", err)
	}
	ret64 := endian.LittleEndian.Uint64(buf.Bytes())
	if ret64 != u64 {
		t.Errorf("mismatch given 0x%x expect 0x%x", ret64, u64)
	}
}

func TestWritePrimitiveBE(t *testing.T) {
	// uint8
	buf := bytes.NewBuffer([]byte{})
	var u8 uint8 = 0xbb
	if err := endian.Write(buf, endian.BigEndian, u8); err != nil {
		t.Errorf("endian.Write error %s", err)
	}
	ret, err := buf.ReadByte()
	if err != nil {
		t.Errorf("ReadByte err=%s", err)
	} else if ret != u8 {
		t.Errorf("mismatch given 0x%x expect 0x%x", ret, u8)
	}

	// uint16
	buf.Reset()
	var u16 uint16 = 0xbbee
	if err := endian.Write(buf, endian.BigEndian, u16); err != nil {
		t.Errorf("endian.Write error %s", err)
	}
	ret16 := endian.BigEndian.Uint16(buf.Bytes())
	if ret16 != u16 {
		t.Errorf("mismatch given 0x%x expect 0x%x", ret16, u16)
	}

	// uint32
	buf.Reset()
	var u32 uint32 = 0xbbeeccff
	if err := endian.Write(buf, endian.BigEndian, u32); err != nil {
		t.Errorf("endian.Write error %s", err)
	}
	ret32 := endian.BigEndian.Uint32(buf.Bytes())
	if ret32 != u32 {
		t.Errorf("mismatch given 0x%x expect 0x%x", ret32, u32)
	}

	// uint64
	buf.Reset()
	var u64 uint64 = 0xbbeeccff00112233
	if err := endian.Write(buf, endian.BigEndian, u64); err != nil {
		t.Errorf("endian.Write error %s", err)
	}
	ret64 := endian.BigEndian.Uint64(buf.Bytes())
	if ret64 != u64 {
		t.Errorf("mismatch given 0x%x expect 0x%x", ret64, u64)
	}
}

func TestWriteByteSlice(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	input := []byte{0xaa, 0xbb, 0xcc, 0xdd}
	if err := endian.Write(buf, endian.LittleEndian, input); err != nil {
		t.Errorf("endian.Write error %s", err)
	}
	ret := buf.Bytes()
	if bytes.Compare(ret, input) != 0 {
		t.Errorf("mismatch\n given:%x\n expect:%x", ret, input)
	}

	// bigendian
	buf.Reset()
	if err := endian.Write(buf, endian.BigEndian, input); err != nil {
		t.Errorf("endian.Write error %s", err)
	}
	ret = buf.Bytes()
	if len(ret) != len(input) {
		t.Errorf("len error given=%d expect=%d", len(ret), len(input))
	}

	length := len(ret)
	expect := make([]byte, length)
	for i := 0; i < length/2; i++ {
		expect[i], expect[length-i-1] = input[length-i-1], input[i]
	}
	if bytes.Compare(expect, ret) != 0 {
		t.Errorf("mismatch\n given=%x expect=%x", ret, expect)
	}

}

func TestWriteByteArray(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	input := [4]byte{0xaa, 0xbb, 0xcc, 0xdd}
	if err := endian.Write(buf, endian.LittleEndian, input); err != nil {
		t.Errorf("endian.Write error %s", err)
	}
	ret := buf.Bytes()
	iexpect := input[:]
	if bytes.Compare(ret, iexpect) != 0 {
		t.Errorf("mismatch\n given:%x\n expect:%x", ret, iexpect)
	}

	// bigendian
	buf.Reset()
	if err := endian.Write(buf, endian.BigEndian, input); err != nil {
		t.Errorf("endian.Write error %s", err)
	}
	ret = buf.Bytes()
	if len(ret) != len(input) {
		t.Errorf("len error given=%d expect=%d", len(ret), len(input))
	}

	length := len(ret)
	expect := make([]byte, length)
	for i := 0; i < length/2; i++ {
		expect[i], expect[length-i-1] = input[length-i-1], input[i]
	}
	if bytes.Compare(expect, ret) != 0 {
		t.Errorf("mismatch\n given=%x expect=%x", ret, expect)
	}

}

func TestWriteStruct(t *testing.T) {
	type S struct {
		B   byte
		U16 uint16
		A   []byte
	}

	s := S{}
	buf := bytes.NewBuffer([]byte{})

	s.B = 0xaa
	s.U16 = 0xbbcc
	s.A = []byte{0xdd, 0xee, 0xff, 0x00, 0x11}

	if err := endian.Write(buf, endian.LittleEndian, s); err != nil {
		t.Errorf("endian.Write err=%s", err)
	}
	expect := []byte{0xaa, 0xcc, 0xbb, 0xdd, 0xee, 0xff, 0x00, 0x11}
	ret := buf.Bytes()
	if bytes.Compare(ret, expect) != 0 {
		t.Errorf("mismatch\n given=%x\n expect=%x", ret, expect)
	}
}

func BenchmarkWriteStruct(b *testing.B) {
	type Sample struct {
		Header   byte
		Reserved [2]byte
		Id       byte
		Data     [4]byte
	}
	s := Sample{}
	br := bytes.NewBuffer([]byte{})
	b.ResetTimer()

	s.Header = 0x7f
	s.Reserved = [2]byte{0xff, 0xff}
	s.Id = 0x51
	s.Data = [4]byte{0xaa, 0xbb, 0xcc, 0xdd}

	for i := 0; i < b.N; i++ {
		br.Reset()
		err := endian.Write(br, endian.LittleEndian, s)
		if err != nil {
			b.Fatalf("error:%s", err)
		}
	}
}

func BenchmarkBinaryWriteStruct(b *testing.B) {
	type Sample struct {
		Header   byte
		Reserved [2]byte
		Id       byte
		Data     [4]byte
	}
	s := Sample{}
	br := bytes.NewBuffer([]byte{})
	b.ResetTimer()

	s.Header = 0x7f
	s.Reserved = [2]byte{0xff, 0xff}
	s.Id = 0x51
	s.Data = [4]byte{0xaa, 0xbb, 0xcc, 0xdd}

	for i := 0; i < b.N; i++ {
		br.Reset()
		err := binary.Write(br, binary.LittleEndian, s)
		if err != nil {
			b.Fatalf("error:%s", err)
		}
	}
}
