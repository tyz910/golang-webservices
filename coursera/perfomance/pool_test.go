package main

import (
	"bytes"
	"encoding/json"
	"sync"
	"testing"
)

const iterNum = 100

type PublicPage struct {
	ID          int
	Name        string
	Url         string
	OwnerID     int
	ImageUrl    string
	Tags        []string
	Description string
	Rules       string
}

var CoolGolangPublic = PublicPage{
	ID:          1,
	Name:        "CoolGolangPublic",
	Url:         "http://example.com",
	OwnerID:     100500,
	ImageUrl:    "http://example.com/img.png",
	Tags:        []string{"programming", "go", "golang"},
	Description: "Best page about golang programming",
	Rules:       "",
}

var Pages = []PublicPage{
	CoolGolangPublic, 
	CoolGolangPublic, 
	CoolGolangPublic,
}

func BenchmarkAllocNew(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			data := bytes.NewBuffer(make([]byte, 0, 64))
			_ = json.NewEncoder(data).Encode(Pages)
		}
	})
}

var dataPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 64))
	},
}

func BenchmarkAllocPool(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			data := dataPool.Get().(*bytes.Buffer)
			_ = json.NewEncoder(data).Encode(Pages)
			data.Reset()
			dataPool.Put(data)
		}
	})
}

/*
	go test -bench . -benchmem pool_test.go
*/
