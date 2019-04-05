package main

import (
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/proto"
)

func main() {
	sess := &Session{
		Login:     "rvasily",
		Useragent: "Chrome",
	}

	dataJson, _ := json.Marshal(sess)

	fmt.Printf("dataJson\nlen %d\n%v\n", len(dataJson), dataJson)

	/*
		40 байт
		{"login":"rvasily","useragent":"Chrome"}
	*/

	dataPb, _ := proto.Marshal(sess)
	fmt.Printf("dataPb\nlen %d\n%v\n", len(dataPb), dataPb)

	/*
		17 байт
		[10 7 114 118 97 115 105 108 121 18 6 67 104 114 111 109 101]

			10 // номер поля + тип
			7  // длина данных
				114 118 97 115 105 108 121
			18 // номер поля + тип
			6  // длина данных
				67 104 114 111 109 101
	*/

}
