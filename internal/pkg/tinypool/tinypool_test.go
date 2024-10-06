package tinypool

import (
	"bytes"
	"strconv"
	"testing"
)

func TestPool(t *testing.T) {
	bufferPool := New[bytes.Buffer](BufReset)
	buff := bufferPool.Get()
	defer bufferPool.Free(buff)
	buff.WriteString("hello,")
	buff.WriteString("gopher,czyt")
	t.Log(buff.String())
}

func BenchmarkPool(b *testing.B) {
	buffReset := func(b *bytes.Buffer) {
		b.Reset()
	}
	bufferPool := New[bytes.Buffer](buffReset)
	for i := 0; i < b.N; i++ {
		buff := bufferPool.Get()
		buff.WriteString(strconv.Itoa(i))
		//b.Log(buff.String())
		bufferPool.Free(buff)
	}

}
