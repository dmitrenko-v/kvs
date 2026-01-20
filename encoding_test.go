package main

import (
	"testing"
)

func TestBytesToIntCorrectPositive(t *testing.T) {
	intBytes := []byte("153")

	want := 153

	res, err := bytesToInt(intBytes)

	if res != want || err != nil {
		t.Errorf("bytesToInt([]byte('153') = %v, expected: %v, err: %v", res, want, err)
	}
}

func TestBytesToIntCorrectNegative(t *testing.T) {
	intBytes := []byte("-153")

	want := -153

	res, err := bytesToInt(intBytes)

	if res != want || err != nil {
		t.Errorf("bytesToInt([]byte('153') = %v, expected: %v, err: %v", res, want, err)
	}
}

func TestBytesToIntNotNumber(t *testing.T) {
	intBytes := []byte("abcdef")

	res, err := bytesToInt(intBytes)

	if res != 0 || err != ErrInvalidIntVal {
		t.Errorf("bytesToInt([]byte('153') = %v, expected: %v, err: %v", res, ErrInvalidIntVal, err)
	}
}

func TestBytesToIntEmptyString(t *testing.T) {
	intBytes := []byte("")

	res, err := bytesToInt(intBytes)

	if res != 0 || err != ErrInvalidIntVal {
		t.Errorf("bytesToInt([]byte('153') = %v, expected: %v, err: %v", res, ErrInvalidIntVal, err)
	}
}

func TestBytesToIntOnlyMinus(t *testing.T) {
	intBytes := []byte("-")

	res, err := bytesToInt(intBytes)

	if res != 0 || err != ErrInvalidIntVal {
		t.Errorf("bytesToInt([]byte('153') = %v, expected: %v, err: %v", res, ErrInvalidIntVal, err)
	}
}

func TestDecodeBoolTrue(t *testing.T) {
	boolBytes := []byte{0x01}

	expected := "true"

	res := decodeBool(boolBytes)

	if res != expected {
		t.Errorf("decodeBool([]byte{0x01}) = %v, expected %v", res, expected)
	}
}

func TestDecodeBoolFalse(t *testing.T) {
	boolBytes := []byte{0x00}

	expected := "false"

	res := decodeBool(boolBytes)

	if res != expected {
		t.Errorf("decodeBool([]byte{0x00}) = %v, expected %v", res, expected)
	}
}

func TestDecodeInt(t *testing.T) {
	bytes1024 := []byte{0x0, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

	expected := "1024"

	res := decodeInt(bytes1024)

	if res != expected {
		t.Errorf("decodeInt(<1024 as byte slice>) = %v, expected %v", res, expected)
	}
}
