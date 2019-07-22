package jsonstream

import (
	"fmt"
	"strings"
	"testing"
)

var testJson string = `{
	"count": 2,
	"crew": [
		{
			"name" : "Spock",
			"age" : 38
		},
		{
			"name" : "Chekov",
			"age" : 26
		},
		{
			"name" : "Scotty",
			"age" : 33
		}
	],
	"captain": {
		"name" : "Kirk",
		"age" : 35
	}
}`

type Person struct {
	Name string `json:"name"`
	Age int `json:"age"`
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func TestDecode(t *testing.T) {
	d := Start(strings.NewReader(testJson))
	it, err := d.GetIterator()
	must(err)

	dict := it.(*DictIterator)
	for dict.Next() {
		key, dec := dict.Value()

		switch key {
		case "count":
			it2, err := dec.GetIterator()
			must(err)

			single := it2.(*SingleIterator)
			count := int(single.Value().(float64))
			fmt.Println("count:", count)
		case "crew":
			it, err := dec.GetIterator()
			must(err)

			arr := it.(*ArrayIterator)
			for arr.Next() {
				index, dec2 := arr.Value()
				var p Person
				must(dec2.Decode(&p))
				fmt.Println("crew:", p)

				//只接受前面2个成员，之后的全部放弃
				if index >= 2 {
					break
				}
			}
			must(arr.Finish())
		case "captain":
			var p Person
			must(dec.Decode(&p))
			fmt.Println("captain:", p)
		}
	}
	must(dict.Finish())
}
