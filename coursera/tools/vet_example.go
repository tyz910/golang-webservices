package main

import (
	"fmt"
)

func test_error(is_ok bool) error {
	if !is_ok {
		fmt.Errorf("failed")
	}
	return nil
}

func main() {
	flag := true
	result := test_error(flag)
	fmt.Printf("result is\n", result)
	fmt.Printf("%v is %v\n", flag)
}
