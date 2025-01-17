package common

import "encoding/hex"

func FromHex(str string) []byte {
	if len(str) < 2 {
		s, _ := hex.DecodeString(str)
		return s
	} else {
		if str[:2] == "0x" || str[:2] == "0X" {
			str = str[2:]
		}
		s, _ := hex.DecodeString(str)
		return s
	}
}
