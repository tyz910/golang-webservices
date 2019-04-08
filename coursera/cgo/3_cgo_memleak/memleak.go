/*
	https://golang.org/cmd/cgo/#hdr-Passing_pointers
*/

package main

// #include <stdio.h>
// #include <stdlib.h>
import "C"
import "unsafe"

func print(s string) {
	cs := C.CString(s) // переход в другую вселенную
	defer C.free(unsafe.Pointer(cs))
	println(cs)
	// СИ-шные не собираются через ГО-шный сборщик мусора, их надо освобождать руками
	// закомментируйте эту строку и запустите программу - начнётся утечка памяти
	// C.free(unsafe.Pointer(cs))
}

func main() {
	for {
		print("Hello World")
	}
}
