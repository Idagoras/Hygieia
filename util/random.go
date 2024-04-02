package util

import "math/rand"

var alphaTable = []byte{
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'a', 'b', 'c', 'd', 'e', 'd', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y',
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(length int) string {
	res := make([]byte, length)
	for i := range res {
		if i == 0 || i == length-1 {
			res[i] = alphaTable[RandomInt(0, int64(len(alphaTable))-1)]
			continue
		}
		res[i] = byte(48 + RandomInt(0, 9))
	}
	return string(res)
}

func GenerateSecretCode(length int) string {
	res := make([]byte, length)
	for i := range res {
		if i == 0 || i == length-1 {
			res[i] = alphaTable[RandomInt(0, int64(len(alphaTable))-1)]
			continue
		}
		res[i] = byte(48 + RandomInt(0, 9))
	}
	return string(res)
}
