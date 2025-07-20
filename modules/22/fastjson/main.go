package main

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fastjson"
)

type Nested struct {
	Key string `json:"key"`
}

func main() {

	var p fastjson.Parser

	jsonData := `{"foo": "bar", "baz": 123, "arr": [1, 2, 3], "nested": {"key": "value"}}`

	v, err := p.Parse(jsonData)
	if err != nil {
		panic(err)
	}

	foo := v.GetStringBytes("foo")
	baz := v.GetInt("baz")
	arr := v.GetArray("arr")
	nested := v.Get("nested", "key")

	fmt.Println("foo:", string(foo))
	fmt.Println("baz:", baz)
	fmt.Println("arr:", arr)
	fmt.Println("nested key:", string(nested.GetStringBytes()))

	nested2 := v.GetObject("nested")
	if nested2 != nil {
		fmt.Println("nested2 key:", nested2.Get("key"))
	} else {
		fmt.Println("nested2 is nil")
	}

	nestedJSON := v.Get("nested").String()
	nestedStruct := Nested{}
	if err := json.Unmarshal([]byte(nestedJSON), &nestedStruct); err != nil {
		panic(err)
	}
	fmt.Println("nested struct key:", nestedStruct.Key)
	fmt.Println("Completed processing JSON data.")
}
