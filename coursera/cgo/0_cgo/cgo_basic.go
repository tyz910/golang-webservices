package main

/*
// Всё что в комментари над import "C" является кодом на C code и будет скомпилирован при помощи GCC.
// У вас должен быть установлен GCC

int Multiply(int a, int b) {
    return a * b;
}
*/
import "C" //это псевдо-пакет, он реализуется компилятором
import "fmt"

func main() {
	a := 2
	b := 3
	// для того чтобы вызвать СИшный крод надо добавить префикс "С."
	// там же туда надо передать СИшные переменные
	res := C.Multiply(C.int(a), C.int(b))
	fmt.Printf("Multiply in C: %d * %d = %d\n", a, b, int(res))
	fmt.Printf("с-var internals %T = %+v\n", res, res)
}
