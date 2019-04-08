package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Config struct {
	Comments bool `json:"comments"`
	Limit    int
	Servers  []string
}

var (
	config = &Config{}
)

func main() {
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatalln("cant read config file:", err)
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		log.Fatalln("cant parse config:", err)
	}

	if config.Comments {
		fmt.Println("Comments per page", config.Limit)
		fmt.Println("Comments services", config.Servers)
	} else {
		fmt.Println("Comments disabled")
	}
}
