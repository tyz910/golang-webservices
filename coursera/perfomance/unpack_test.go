package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"
)

/*
	go test -bench . unpack_test.go
	go test -bench . -benchmem unpack_test.go
*/

var (
	data = []byte{
		128, 36, 17, 0,

		9, 0, 0, 0,
		118, 46, 114, 111, 109, 97, 110, 111, 118,

		16, 0, 0, 0,
	}
)

type User struct {
	ID       int
	RealName string `cgen:"-"`
	Login    string
	Flags    int
}

func BenchmarkCodegen(b *testing.B) {
	u := &User{}
	for i := 0; i < b.N; i++ {
		u = &User{}
		u.UnpackBin(data)
	}
}

func BenchmarkReflect(b *testing.B) {
	u := &User{}
	for i := 0; i < b.N; i++ {
		u = &User{}
		UnpackReflect(u, data)
	}
}

func (in *User) UnpackBin(data []byte) error {
	r := bytes.NewReader(data)

	// ID
	var IDRaw uint32
	binary.Read(r, binary.LittleEndian, &IDRaw)
	in.ID = int(IDRaw)

	// Login
	var LoginLenRaw uint32
	binary.Read(r, binary.LittleEndian, &LoginLenRaw)
	LoginRaw := make([]byte, LoginLenRaw)
	binary.Read(r, binary.LittleEndian, &LoginRaw)
	in.Login = string(LoginRaw)

	// Flags
	var FlagsRaw uint32
	binary.Read(r, binary.LittleEndian, &FlagsRaw)
	in.Flags = int(FlagsRaw)
	return nil
}

func UnpackReflect(u interface{}, data []byte) error {
	r := bytes.NewReader(data)

	val := reflect.ValueOf(u).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		if typeField.Tag.Get("unpack") == "-" {
			continue
		}

		switch typeField.Type.Kind() {
		case reflect.Int:
			var value int
			binary.Read(r, binary.LittleEndian, &value)
			valueField.Set(reflect.ValueOf(value))
		case reflect.String:
			var lenRaw int
			binary.Read(r, binary.LittleEndian, &lenRaw)

			dataRaw := make([]byte, lenRaw)
			binary.Read(r, binary.LittleEndian, &dataRaw)

			valueField.SetString(string(dataRaw))
		default:
			return fmt.Errorf("bad type: %v for field %v", typeField.Type.Kind(), typeField.Name)
		}
	}

	return nil
}

/*
	go test -bench . unpack_test.go
	go test -bench . -benchmem unpack_test.go

	go test -bench . -benchmem -cpuprofile=cpu.out -memprofile=mem.out -memprofilerate=1 unpack_test.go

	go tool pprof main.test.exe cpu.out
	go tool pprof main.test.exe mem.out

	go get github.com/uber/go-torch
	go-torch main.test.exe cpu.out

*/
