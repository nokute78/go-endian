# go-endian

A library to read/write n-byte big/little endian data.

## Feature

* n-byte endian 
* mixed endian

## Installation

```
$ go get github.com/nokute78/go-endian
```

## Usage

The package supports `binary.Read` like API.

It is an example to decode GUID.
```go
package main

import (
	"bytes"
	"fmt"
	"github.com/nokute78/go-endian"
)

func main() {
	type GUID struct {
		G0 uint32
		G1 uint16
		G2 uint16
		G3 uint16 `endian:"BE"`
		G4 [6]byte
	}

	// EFI System Partition GUID
	// C12A7328-F81F-11D2-BA4B-00A0C93EC93B
	raw := []byte{0x28, 0x73, 0x2a, 0xc1, 0x1f, 0xf8, 0xd2, 0x11, 0xba, 0x4b, 0x00, 0xa0, 0xc9, 0x3e, 0xc9, 0x3b}

	buf := bytes.NewReader(raw)
	guid := GUID{}
	endian.Read(buf, endian.LittleEndian, &guid)

	fmt.Printf("GUID:%x-%x-%x-%x-%x\n", guid.G0, guid.G1, guid.G2, guid.G3, guid.G4)
}
```

## Struct Tag

The package supports struct tags.

|Tag|Description|
|---|-----------|
|`` `endian:"skip"` ``|Ignore the field. Offset is updated by the size of the field. It is useful for reserved field.|
|`` `endian:"-"` `` |Ignore the field. Offset is not updated.|
|`` `endian:"BE"` ``|Decode the field as big endian. It is useful for mixed endian data.|
|`` `endian:"LE"` ``|Decode the field as little endian. It is useful for mixed endian data.|


## Document


## License

[Apache License v2.0](https://www.apache.org/licenses/LICENSE-2.0)