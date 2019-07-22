package jsonstream

type SingleIterator struct {
	v interface{}
}

func newSingleIterator(v interface{}) *SingleIterator {
	return &SingleIterator{
		v: v,
	}
}

func (si *SingleIterator) Finish() error {
	return nil
}

func (si *SingleIterator) Next() bool {
	return false
}

func (si *SingleIterator) Value() interface{} {
	return si.v
}
