package bytesEx

import "bytes"

func IsLastByteNewline(data []byte) bool {
	return len(data) > 0 && (bytes.HasSuffix(data, []byte{'\n'}) || bytes.HasSuffix(data, []byte{'\r'}))
}
