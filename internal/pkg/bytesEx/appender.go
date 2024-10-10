package bytesEx

import "bytes"

func AppendPerLine(bytesToAppend []byte, appendContent string) []byte {
	parts := bytes.Split(bytesToAppend, []byte{'\n'})
	var buffer bytes.Buffer
	for i, part := range parts {
		buffer.Write(part)
		buffer.WriteString(appendContent)
		if i < len(parts)-1 {
			buffer.WriteByte('\n')
		}
	}
	return buffer.Bytes()
}
