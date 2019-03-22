// go build gen/* && ./codegen.exe pack/packer.go  pack/marshaller.go
package main

import "fmt"

// lets generate code for this struct
// cgen: binpack
type User struct {
	ID       int
	RealName string `cgen:"-"`
	Login    string
	Flags    int
}

type Avatar struct {
	ID  int
	Url string
}

var test = 42

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

	u := User{}
	u.Unpack(data)
	fmt.Printf("Unpacked user %#v", u)
}
