package main

import (
	"fmt"
)

/*
	go run -ldflags="-X 'main.Version=$(git rev-parse HEAD)' -X 'main.Branch=$(git rev-parse --abbrev-ref HEAD)'" ldflags.go
*/

var (
	Version = ""
	Branch  = ""
)

func main() {
	fmt.Println("[start] starting version ", Version, Branch)
}
