// https://github.com/tschottdorf/goplay/tree/master/cgobench
// go test -bench . -gcflags '-l'    # отключаем инлайнинг
package main

import (
	"testing"
)

func BenchmarkCGO(b *testing.B) {
	CallCgo(b.N) // call `C.(void f() {})` b.N times
}

// BenchmarkGo должен быть вызывать с опцией `-gcflags -l` чтобы отключить инлайнинг
func BenchmarkGo(b *testing.B) {
	CallGo(b.N) // call `func() {}` b.N times
}
