package main

import (
	"encoding/json"
	"testing"
)

var (
	data = []byte(`{"RealName":"Vasily", "Login":"v.romanov", "Status":1, "Flags": 1}`)
	u    = User{}
	c    = Client{}
)

// go test -v -bench=. -benchmem json/*.go
// go test -v -bench=. json/*.go

func BenchmarkDecodeStandart(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = json.Unmarshal(data, &c)
	}
}

func BenchmarkDecodeEasyjson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = u.UnmarshalJSON(data)
	}
}

func BenchmarkEncodeStandart(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(&c)
	}
}

func BenchmarkEncodeEasyjson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = u.MarshalJSON()
	}
}
