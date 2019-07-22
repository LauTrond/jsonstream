package jsonstream

import (
	"encoding/json"
	"fmt"
	"path"
)

type DictIterator struct {
	parent *Decoder
	done bool

	key string
	lastChild *Decoder
	err error
}

func newDictIterator(d *Decoder) *DictIterator {
	return &DictIterator{
		parent: d,
	}
}

func (it *DictIterator) Finish() error {
	for it.Next() {}
	return it.err
}

func (it *DictIterator) Value() (string, *Decoder) {
	return it.key, it.lastChild
}

func (it *DictIterator) Next() bool {
	if it.lastChild != nil {
		it.lastChild.consumeAll()
		it.lastChild = nil
	}
	it.key = ""

	if it.done {
		return false
	}

	token, err := it.parent.jsonDecoder.Token()
	if err != nil {
		it.err = err
		return false
	}

	if d,ok := token.(json.Delim); ok  && d == '}'{
		it.done = true
		return false
	}

	key, ok := token.(string)
	if !ok {
		it.err = fmt.Errorf("%s DictIterator: unexpected token (%T)%v",
			it.parent.jsonPath, token, token)
	}

	it.key = key
	it.lastChild = newDecoder(it.parent.jsonDecoder, path.Join(it.parent.jsonPath, key))
	return true
}
