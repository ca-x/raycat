package bytesEx

import "testing"

func TestIsLastByteNewline(t *testing.T) {
	data1 := []byte("Hello, World\n")
	isLastByteNewline := IsLastByteNewline(data1)
	t.Log("isLastByteNewline:", isLastByteNewline)
}
