package main

import (
	"encoding/json"
	"reflect"
	"testing"
	// "fmt"
)

type Simple struct {
	ID       int
	Username string
	Active   bool
}

type IDBlock struct {
	ID int
}

func TestSimple(t *testing.T) {
	expected := &Simple{
		ID:       42,
		Username: "rvasily",
		Active:   true,
	}
	jsonRaw, _ := json.Marshal(expected)
	// fmt.Println(string(jsonRaw))

	var tmpData interface{}
	json.Unmarshal(jsonRaw, &tmpData)

	result := new(Simple)
	err := i2s(tmpData, result)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("results not match\nGot:\n%#v\nExpected:\n%#v", result, expected)
	}
}

type Complex struct {
	SubSimple  Simple
	ManySimple []Simple
	Blocks     []IDBlock
}

func TestComplex(t *testing.T) {
	smpl := Simple{
		ID:       42,
		Username: "rvasily",
		Active:   true,
	}
	expected := &Complex{
		SubSimple:  smpl,
		ManySimple: []Simple{smpl, smpl},
		Blocks:     []IDBlock{IDBlock{42}, IDBlock{42}},
	}

	jsonRaw, _ := json.Marshal(expected)
	// fmt.Println(string(jsonRaw))

	var tmpData interface{}
	json.Unmarshal(jsonRaw, &tmpData)

	result := new(Complex)
	err := i2s(tmpData, result)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("results not match\nGot:\n%#v\nExpected:\n%#v", result, expected)
	}
}

func TestSlice(t *testing.T) {
	smpl := Simple{
		ID:       42,
		Username: "rvasily",
		Active:   true,
	}
	expected := []Simple{smpl, smpl}

	jsonRaw, _ := json.Marshal(expected)

	var tmpData interface{}
	json.Unmarshal(jsonRaw, &tmpData)

	result := []Simple{}
	err := i2s(tmpData, &result)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("results not match\nGot:\n%#v\nExpected:\n%#v", result, expected)
	}
}

type ErrorCase struct {
	Result   interface{}
	JsonData string
}

// аккуратно в этом тесте
// писать надо именно в то что пришло
func TestErrors(t *testing.T) {
	cases := []ErrorCase{
		// "Active":"DA" - string вместо bool
		ErrorCase{
			&Simple{},
			`{"ID":42,"Username":"rvasily","Active":"DA"}`,
		},
		// "ID":"42" - string вместо int
		ErrorCase{
			&Simple{},
			`{"ID":"42","Username":"rvasily","Active":true}`,
		},
		// "Username":100500 - int вместо string
		ErrorCase{
			&Simple{},
			`{"ID":42,"Username":100500,"Active":true}`,
		},
		// "ManySimple":{} - ждём слайс, получаем структуру
		ErrorCase{
			&Complex{},
			`{"SubSimple":{"ID":42,"Username":"rvasily","Active":true},"ManySimple":{}}`,
		},
		// "SubSimple":true - ждём структуру, получаем bool
		ErrorCase{
			&Complex{},
			`{"SubSimple":true,"ManySimple":[{"ID":42,"Username":"rvasily","Active":true}]}`,
		},
		// ожидаем структуру - пришел массив
		ErrorCase{
			&Simple{},
			`[{"ID":42,"Username":"rvasily","Active":true}]`,
		},
		// Simple{} ( без амперсанта, т.е. структура, а не указатель на структуру )
		// пришел не ссылочный тип - мы не сможем вернуть результат
		ErrorCase{
			Simple{},
			`{"ID":42,"Username":"rvasily","Active":true}`,
		},
	}
	for idx, item := range cases {
		var tmpData interface{}
		json.Unmarshal([]byte(item.JsonData), &tmpData)
		inType := reflect.ValueOf(item.Result).Type()
		err := i2s(tmpData, item.Result)
		outType := reflect.ValueOf(item.Result).Type()

		if err == nil {
			t.Errorf("[%d] expected error here", idx)
			continue
		}
		if inType != outType {
			t.Errorf("results type not match\nGot:\n%#v\nExpected:\n%#v", outType, inType)
		}
	}
}
