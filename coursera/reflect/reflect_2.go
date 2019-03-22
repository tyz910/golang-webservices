package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

type User struct {
	ID       int
	RealName string `unpack:"-"`
	Login    string
	Flags    int
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
			var value uint32
			binary.Read(r, binary.LittleEndian, &value)
			valueField.Set(reflect.ValueOf(int(value)))
		case reflect.String:
			var lenRaw uint32
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

func main() {
	/*
		perl -E '$b = pack("L L/a* L", 1_123_456, "v.romanov", 16);
			print map { ord.", "  } split("", $b); '
	*/
	data := []byte{
		128, 36, 17, 0,

		9, 0, 0, 0,
		118, 46, 114, 111, 109, 97, 110, 111, 118,

		16, 0, 0, 0,
	}
	u := new(User)
	err := UnpackReflect(u, data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", u)
}
