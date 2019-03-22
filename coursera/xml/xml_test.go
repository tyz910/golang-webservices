package main

import (
	"testing"
)

func BenchmarkCountStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CountStruct()
	}
}

func BenchmarkCountDecoder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CountDecoder()
	}
}
