package scanner

import "bytes"

func match(base, pattern []byte, index int) bool {
	length := len(pattern)
	if index+length >= len(base) {
		return false
	}
	return bytes.Equal(bytes.ToLower(base[index:index+length]), pattern)
}
