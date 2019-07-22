//package jsonstream 提供简单的JSON流式处理：
//所谓流式处理，就是可以在不把整个JSON数据加载到内存的前提下完成遍历。
//package jsonstream是package json的封装，
//当你需要JSON流式处理，jsonstream比使用标准库json.Decoder更简单。
package jsonstream

import (
	"encoding/json"
	"fmt"
	"io"
)

var ErrDecoderUsed = fmt.Errorf("decoder used")

func Start(r io.Reader) *Decoder {
	return newDecoder(json.NewDecoder(r), "")
}

func newDecoder(d *json.Decoder, jsonPath string) *Decoder {
	return &Decoder{
		jsonDecoder: d,
		jsonPath: jsonPath,
	}
}

//每个Decoder对应一个JSON对象，可能它可能是：
//数或字符串
//数组：[...]
//结构体：{...}
type Decoder struct{
	jsonDecoder *json.Decoder
	jsonPath string
	iterator Iterator
	used bool
}

//把整个JSON对象解码到指针v指向的变量，等同json.Decoder.Decode
func (d *Decoder) Decode(v interface{}) error {
	if d.used {
		return ErrDecoderUsed
	}
	d.used = true

	err := d.jsonDecoder.Decode(v)
	if err != nil {
		err = fmt.Errorf("%s: %s", d.jsonPath, err.Error())
	}
	return err
}

//分解JSON对象，获得一个迭代器，迭代获取子对象。
//返回迭代器的类型可能是：*SingleIterator, *DictIterator, *ArrayIterator
func (d *Decoder) GetIterator() (Iterator,error) {
	if d.used {
		return nil, ErrDecoderUsed
	}
	d.used = true

	t, err := d.jsonDecoder.Token()
	if err != nil {
		return nil, err
	}

	var iterator Iterator

	switch v := t.(type) {
	case nil,bool,float64,json.Number,string:
		iterator = newSingleIterator(v)
	case json.Delim:
		switch v {
		case '{':
			iterator = newDictIterator(d)
		case '[':
			iterator = newArrayIterator(d)
		default:
			panic(fmt.Errorf("unexpected json.Delim '%c'", v))
		}
	default:
		//json.Token 不包含的类型
		panic(fmt.Errorf("unknown json.Token %T", v))
	}

	d.iterator = iterator
	return iterator, nil
}

func (d *Decoder) consumeAll() error {
	if d.iterator == nil {
		if d.used {
			//已经执行Decode()
			return nil
		}
		_, err := d.GetIterator()
		if err != nil {
			return err
		}
	}
	return d.iterator.Finish()
}

type Iterator interface{
	Next() bool
	Finish() error
}
