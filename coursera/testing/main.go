package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID int
}

var data = map[string][]byte{
	"ok":   []byte(`{"ID": 27}`),
	"fail": []byte(`{"ID": 27`),
}

func GetUser(key string) (*User, error) {
	if jsonStr, ok := data[key]; ok {
		user := &User{}
		err := json.Unmarshal(jsonStr, user)
		if err != nil {
			return nil, fmt.Errorf("Cant decode json")
		}
		return user, nil
	}
	return nil, fmt.Errorf("User doesnt exist")
}
