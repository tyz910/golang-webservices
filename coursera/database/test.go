package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

type RSS struct {
	Items []Item `xml:"channel>item"`
}

type Item struct {
	URL   string `xml:"guid"`
	Title string `xml:"title"`
}

func getData() interface{} {
	// return Item{"1", "1"}
	return &RSS{
		Items: []Item{
			Item{"1", "1"},
			Item{"2", "2"},
			Item{"3", "3"},
			Item{"4", "4"},
			Item{"5", "5"},
		},
	}
}

type CacheItem struct {
	Data json.RawMessage
	// Tags map[string]int
}

type CacheItemStore struct {
	Data interface{}
}

var data []byte

func getDeep(in interface{}) {
	in = getData()
}

func unpackData(in interface{}, mode bool) error {

	if mode == false {
		fmt.Println("set in")

		tmpIn := getData()

		fmt.Println("types", reflect.TypeOf(tmpIn), reflect.TypeOf(in))
		fmt.Println("types equal", reflect.TypeOf(tmpIn) == reflect.TypeOf(in))

		// fmt.Println(reflect.ValueOf(in))
		// reflect.ValueOf(in).Set(tmpIn)

		// *tt := *tmpIn

		rvp := reflect.ValueOf(in)
		fmt.Printf("rvp %+v %+v\n", rvp, rvp.Kind())
		if rvp.Kind() != reflect.Ptr {
			panic("Ожидается указатель")
		}

		rvp2 := reflect.ValueOf(tmpIn)
		fmt.Printf("rvp2 %+v %+v\n", rvp2, rvp2.Kind())
		if rvp2.Kind() != reflect.Ptr {
			panic("Ожидается указатель")
		}

		// rvp.Set(rvp2)

		rv := reflect.Indirect(rvp)
		fmt.Printf("rv %+v %+v\n", rv, rv.Kind())
		rv.Set(reflect.Indirect(rvp2))

		// *in = tmpIn

		fmt.Println("in in func", in)
		// getDeep(in)
		return nil
	}

	d := CacheItemStore{
		Data: getData(),
	}

	d1, err := json.Marshal(d)
	fmt.Println(err)

	d2 := CacheItem{}

	err = json.Unmarshal(d1, &d2)
	// fmt.Println(d2, err)

	err = json.Unmarshal(d2.Data, &in)
	// fmt.Println(in, err)

	return nil
}

func main() {

	fmt.Println(fmt.Sprint(time.Now().Unix()))

	var i = int(time.Now().Unix())
	fmt.Println(i)

	// d1, _ := json.Marshal(getData())

	// fmt.Println(d2)

	// d3 := RSS{}
	// err := unpackData(&d3, true)
	// fmt.Println("result", d3, err)

	d3 := &RSS{}
	// var v interface{} = d3
	err := unpackData(d3, false)
	// fmt.Println("v", v)
	fmt.Println("d3", d3, err)

	time.Sleep(10 * time.Second)

	fmt.Println("d3", d3, err)

}
