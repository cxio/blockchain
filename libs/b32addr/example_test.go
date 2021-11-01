// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Keep in sync with ../base64/example_test.go.

package b32addr_test

import (
	"fmt"
	"os"

	"github.com/cxio/cxsuite/bcutil/b32addr"
)

func ExampleEncoding_EncodeToString() {
	data := []byte("any + old & data")
	str := b32addr.StdEncoding.EncodeToString(data)
	fmt.Println(str)
	// Output:
	// MFXHSIBLEBXWYZBAEYQGIYLUME======
}

func ExampleEncoding_DecodeString() {
	str := "ONXW2ZJAMRQXIYJAO5UXI2BAAAQGC3TEEDX3XPY="
	data, err := b32addr.StdEncoding.DecodeString(str)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("%q\n", data)
	// Output:
	// "some data with \x00 and \ufeff"
}

func ExampleNewEncoder() {
	input := []byte("foo\x00bar")
	encoder := b32addr.NewEncoder(b32addr.StdEncoding, os.Stdout)
	encoder.Write(input)
	// Must close the encoder when finished to flush any partial blocks.
	// If you comment out the following line, the last partial block "r"
	// won't be encoded.
	encoder.Close()
	// Output:
	// MZXW6ADCMFZA====
}
