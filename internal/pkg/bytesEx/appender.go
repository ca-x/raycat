package bytesEx

import "bytes"

func AppendPerLine(bytesToAppend []byte, appendContent string) []byte {
	parts := bytes.Split(bytesToAppend, []byte{'\n'})
	if len(parts) == 1 {
		return bytesToAppend
	}
	var buffer bytes.Buffer
	for i, part := range parts {
		buffer.Write(part)
		// only handle contains # ones for this moment
		if bytes.Contains(part, []byte("#")) {
			buffer.WriteString(appendContent)
		}
		if i < len(parts)-1 {
			buffer.WriteByte('\n')
		}
	}
	return buffer.Bytes()
}
