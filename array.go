package jsonstream

import (
	"encoding/json"
	"fmt"
)

type ArrayIterator struct {
	parent *Decoder
	done bool

	index int
	lastChild *Decoder
	err error
}

func newArrayIterator(d *Decoder) *ArrayIterator {
	return &ArrayIterator{
		parent: d,
	}
}

func (it *ArrayIterator) Finish() error {
	for it.Next() {}
	return it.err
}

func (it *ArrayIterator) Value() (int, *Decoder) {
	return it.index, it.lastChild
}

func (it *ArrayIterator) Next() bool {
	if it.lastChild != nil {
		it.lastChild.consumeAll()
		it.lastChild = nil
	}

	if it.done {
		return false
	}

	if !it.parent.jsonDecoder.More() {
		token, err := it.parent.jsonDecoder.Token()
		if err != nil {
			it.err = err
			return false
		}
		if d,ok := token.(json.Delim); ok && d == ']' {
			it.done = true
			return false
		}
		it.err = fmt.Errorf("%s DictIterator: unexpected token (%T)%v",
			it.parent.jsonPath, token, token)
		return false
	}

	it.index++
	it.lastChild = newDecoder(it.parent.jsonDecoder,
		fmt.Sprintf("%s[%d]",it.parent.jsonPath,it.index))
	return true
}
