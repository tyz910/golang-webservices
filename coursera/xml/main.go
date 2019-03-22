package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

type User struct {
	ID      int    `xml:"id,attr"`
	Login   string `xml:"login"`
	Name    string `xml:"name"`
	Browser string `xml:"browser"`
}

type Users struct {
	Version string `xml:"version,attr"`
	List    []User `xml:"user"`
}

func CountStruct() {
	logins := make([]string, 0)
	v := new(Users)
	err := xml.Unmarshal(xmlData, &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	for _, u := range v.List {
		logins = append(logins, u.Login)
	}
}

func CountDecoder() {
	input := bytes.NewReader(xmlData)
	decoder := xml.NewDecoder(input)
	logins := make([]string, 0)
	var login string
	for {
		tok, tokenErr := decoder.Token()
		if tokenErr != nil && tokenErr != io.EOF {
			fmt.Println("error happend", tokenErr)
			break
		} else if tokenErr == io.EOF {
			break
		}
		if tok == nil {
			fmt.Println("t is nil break")
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			if tok.Name.Local == "login" {
				if err := decoder.DecodeElement(&login, &tok); err != nil {
					fmt.Println("error happend", err)
				}
				logins = append(logins, login)
			}
		}
	}
}

/*
	go test -bench . -benchmem xml_test.go
*/

func main() {
	CountStruct()
	CountDecoder()
}

var xmlData = []byte(`<?xml version="1.0" encoding="utf-8"?>
	<users>
		<user id="1">
			<login>user1</login>
			<name>Василий Романов</name>
			<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36
	</browser>
		</user>
		<user id="2">
			<login>user2</login>
			<name>Иван Иванов</name>
			<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36
	</browser>
		</user>
		<user id="2">
			<login>user3</login>
			<name>Иван Петров</name>
			<browser>Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0; Trident/5.0)</browser>
		</user>
		<user id="1">
			<login>user1</login>
			<name>Василий Романов</name>
			<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36
	</browser>
		</user>
		<user id="2">
			<login>user2</login>
			<name>Иван Иванов</name>
			<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36
	</browser>
		</user>
		<user id="2">
			<login>user3</login>
			<name>Иван Петров</name>
			<browser>Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0; Trident/5.0)</browser>
		</user>
		<user id="2">
			<login>user3</login>
			<name>Иван Петров</name>
			<browser>Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0; Trident/5.0)</browser>
		</user>
		<user id="1">
			<login>user1</login>
			<name>Василий Романов</name>
			<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36
	</browser>
		</user>
		<user id="2">
			<login>user2</login>
			<name>Иван Иванов</name>
			<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36
	</browser>
		</user>
		<user id="2">
			<login>user3</login>
			<name>Иван Петров</name>
			<browser>Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0; Trident/5.0)</browser>
		</user>
		<user id="2">
			<login>user3</login>
			<name>Иван Петров</name>
			<browser>Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0; Trident/5.0)</browser>
		</user>
		<user id="1">
			<login>user1</login>
			<name>Василий Романов</name>
			<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36
	</browser>
		</user>
		<user id="2">
			<login>user2</login>
			<name>Иван Иванов</name>
			<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36
	</browser>
		</user>
		<user id="2">
			<login>user3</login>
			<name>Иван Петров</name>
			<browser>Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0; Trident/5.0)</browser>
		</user>
	</users>`)
