package tinypool

import "bytes"

var (
	BufReset = func(b *bytes.Buffer) {
		b.Reset()
	}
)
