package b32addr_test

import (
	"testing"

	"github.com/cxio/cxsuite/b32addr"
)

func BenchmarkEncodeToString(b *testing.B) {
	data := make([]byte, 50)
	b.SetBytes(int64(len(data)))
	for i := 0; i < b.N; i++ {
		b32addr.StdEncoding.EncodeToString(data)
	}
}

func BenchmarkDecodeString(b *testing.B) {
	data := b32addr.StdEncoding.EncodeToString(make([]byte, 50))
	b.SetBytes(int64(len(data)))
	for i := 0; i < b.N; i++ {
		b32addr.StdEncoding.DecodeString(data)
	}
}
