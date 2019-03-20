package main

import "fmt"

func getSomeVars() string {
	fmt.Println("getSomeVars execution")
	return "getSomeVars result"
}

func main() {
	defer fmt.Println("After work")
	defer func() {
		fmt.Println(getSomeVars())
	}()
	fmt.Println("Some userful work")
}
