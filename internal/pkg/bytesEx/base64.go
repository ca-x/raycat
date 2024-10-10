package bytesEx

func IsBase64(data []byte) bool {
	if len(data)%4 != 0 {
		return false
	}
	for _, b := range data {
		if !((b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') || b == '+' || b == '/' || b == '=') {
			return false
		}
	}
	return true
}
