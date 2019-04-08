// вызов cgo не бесплатный
// https://github.com/tschottdorf/goplay/tree/master/cgobench
// go test -bench . -gcflags '-l'    # disable inlining for fairness

package main

//#include <unistd.h>
//void foo() { }
//void fooSleep() { sleep(100); }
import "C"

func foo() {}

func CallCgo(n int) {
	for i := 0; i < n; i++ {
		C.foo()
	}
}

func CallGo(n int) {
	for i := 0; i < n; i++ {
		foo()
	}
}
