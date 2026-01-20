package main

import (
	"encoding/binary"
	"strconv"
	"unicode"
)

func bytesToInt(intBytes []byte) (res int, err error) {
	if len(intBytes) < 1 {
		return 0, ErrInvalidIntVal
	}

	negative := intBytes[0] == '-'

	if negative && len(intBytes) < 2 {
		return 0, ErrInvalidIntVal
	}

	var firstDigitIx int
	if negative {
		firstDigitIx = 1
	} else {
		firstDigitIx = 0
	}

	factor := 1
	numLen := len(intBytes)

	for i := numLen - 1; i >= firstDigitIx; i-- {
		ch := rune(intBytes[i])
		if !unicode.IsDigit(ch) {
			return 0, ErrInvalidIntVal
		}

		res += int(ch-'0') * factor

		factor *= 10
	}

	if negative {
		res = -res
	}

	return res, nil
}

// the point of decoding is to translate internal byte representation of data into strings
// for them to be sent in response
func decodeBool(boolBytesVal []byte) string {
	if boolBytesVal[0] == 0x01 {
		return "true"
	}

	return "false"
}

func decodeInt(intBytesVal []byte) string {
	intVal := int(binary.NativeEndian.Uint64(intBytesVal))

	return strconv.Itoa(intVal)
}
