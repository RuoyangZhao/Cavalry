package bufpool

import (
	"math/bits"
	"sync"
)

const (
	numClasses = 7
	classMulti = 2
	shift      = 10
)

var (
	bufpools [numClasses]sync.Pool
	classes  [numClasses]int
)

func create(class int) func() any { return func() any { return make([]byte, class) } }

func init() {
	var baseline int = 1 << shift
	for index := range numClasses {
		class := baseline
		bufpools[index] = sync.Pool{New: create(class)}
		classes[index] = class

		baseline *= classMulti
	}
}

func indexOf(size int) int {
	if size <= 0 {
		return -1
	}

	var index int = bits.Len(uint((size - 1))) - shift
	if index < 0 {
		return 0 // size is smaller than the smallest class always belongs to the first class
	}
	if index >= numClasses {
		return -1
	}

	return index
}

func ClassOf(size int) int {
	if index := indexOf(size); index >= 0 {
		return classes[index]
	}

	return -1
}

func alloc(size int) ([]byte, bool) {
	if size <= 0 {
		return nil, false
	}

	if index := indexOf(size); index >= 0 {
		buf := bufpools[index].Get().([]byte)

		return buf[:size], false
	}

	return make([]byte, size), true // fallback to heap allocation, and set the flag to true
}

func Alloc(size int) []byte {
	buf, _ := alloc(size) // never mind where the buffer comes from, just return it

	return buf
}

func Free(buf []byte) {
	if buf == nil {
		return
	}

	var capacity int = cap(buf)
	if capacity == 0 {
		return
	}

	if index := indexOf(capacity); index >= 0 {
		clear(buf[:capacity]) // clear buf before putting back to pool
		bufpools[index].Put(buf[:0])
	}
}
