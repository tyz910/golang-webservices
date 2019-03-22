package main

import (
	"testing"
	"reflect"
)

type TestCase struct{
	Key string
	User *User
	IsError bool
}

func TestGetUser(t *testing.T) {
	cases := []TestCase{
		TestCase{"ok", &User{ID: 27}, false},
		// TestCase{"fail", nil, true},
		TestCase{"not_exist", nil, true},
	}

	for caseNum, item := range cases {
		u, err := GetUser(item.Key)

		if item.IsError && err == nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !item.IsError && err != nil {
			t.Errorf("[%d] unexpected error", caseNum, err)
		}
		if !reflect.DeepEqual(u, item.User) {
			t.Errorf("[%d] wrong results: got %+v, expected %+v",
			caseNum, u, item.User)
		}
	}

}

/*
	go test -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html

*/