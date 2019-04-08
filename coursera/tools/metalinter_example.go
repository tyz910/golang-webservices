package main

import (
	"fmt"
)

type MyStruct struct {
	user_id  int
	DataJson []byte
}

func Test_error(is_ok bool) error {
	if !is_ok {
		fmt.Errorf("failed")
	}
	return nil
}

func Test() {
	flag := true
	result := Test_error(flag)
	fmt.Printf("result is\n", result)
	fmt.Printf("%v is %v", flag)
}
