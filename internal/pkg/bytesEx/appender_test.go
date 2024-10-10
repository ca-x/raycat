package bytesEx

import (
	"encoding/base64"
	"testing"
)

var data = "xxxx"

func TestAppendPerLine(t *testing.T) {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(dst, []byte(data))
	if err != nil {
		t.Fatal(err)
	}
	dst = dst[:n]
	line := AppendPerLine(dst, "🤭fucking high")
	s := base64.StdEncoding.EncodeToString(line)
	t.Log(s)
}
