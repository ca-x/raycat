package bytesEx

import (
	"encoding/base64"
	"testing"
)

func TestIsBase64(t *testing.T) {
	origin := "hi,gopher"
	b64 := base64.StdEncoding.EncodeToString([]byte(origin))
	isb64 := IsBase64([]byte(b64))
	t.Log(isb64)
	isB642 := IsBase64([]byte(origin))
	t.Log(isB642)
}
