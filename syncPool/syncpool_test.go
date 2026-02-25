package main

import (
	"bytes"
	"sync"
	"testing"
)

func BenchmarkWithoutPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := new(bytes.Buffer)
		buf.WriteString("test")
		_ = buf.Bytes()
	}
}

func BenchmarkWithPool(b *testing.B) {
	pool := sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}
	for i := 0; i < b.N; i++ {
		buf := pool.Get().(*bytes.Buffer)
		buf.Reset()
		buf.WriteString("test")
		_ = buf.Bytes()
		pool.Put(buf)
	}
}
