package main

import (
	"bufio"
	"io"
	"slices"
	"strings"
	"testing"
)

func initMockReader(input string) *respReader {
	sReader := strings.NewReader(input)
	bReader := bufio.NewReader(sReader)
	respReader := &respReader{reader: bReader}

	return respReader
}

// ================================ readRespLine ======================================
func TestReadRespLineCorrect(t *testing.T) {
	r := initMockReader("*2\r\n")

	expected := []byte("*2")

	res, err := r.readRespLine()

	if !slices.Equal(res, expected) || err != nil {
		t.Errorf("readLine(*2\\r\\n) = %s, expected: %s, err: %v", res, expected, err)
	}
}

func TestReadRespLineEOF(t *testing.T) {
	r := initMockReader("")

	res, err := r.readRespLine()

	if res != nil || err != io.EOF {
		t.Errorf("readLine(\"\") = %v, expected: %v, err: %v", res, io.EOF, err)
	}
}

func TestReadRespLineWithoutCR(t *testing.T) {
	r := initMockReader("*2\n")

	res, err := r.readRespLine()

	if res != nil || err != ErrInvalidRESP {
		t.Errorf("readLine(*2\\n) = %v, expected: %s, err: %v", res, ErrInvalidRESP, err)
	}
}

// ================================ readArray ========================================
func TestReadArrayCorrectArray(t *testing.T) {
	array := []byte("*2")

	err := readArray(array)

	if err != nil {
		t.Errorf("readArray([]byte('*2') errors with %v, expected: %v", err, nil)
	}
}

func TestReadArrayNoArraySymbol(t *testing.T) {
	array := []byte("2")

	err := readArray(array)

	if err != ErrCmdNotArray {
		t.Errorf("readArray([]byte('2') errors with %v, expected: %v", err, ErrCmdNotArray)
	}
}

func TestReadArrayNegativeCount(t *testing.T) {
	array := []byte("*-2")

	err := readArray(array)

	if err != ErrIncorrectDataLen {
		t.Errorf("readArray([]byte('*-2') errors with %v, expected: %v", err, ErrIncorrectDataLen)
	}
}

func TestReadArrayCountNotNumeric(t *testing.T) {
	array := []byte("*asdf")

	err := readArray(array)

	if err != ErrInvalidIntVal {
		t.Errorf("readArray([]byte('*asdf') errors with %v, expected: %v", err, ErrInvalidIntVal)
	}
}

// ================================ readBulkString ========================================
func TestReadBulkStringCorrect(t *testing.T) {
	r := initMockReader("2\r\nHI\r\n")

	expected := []byte("HI")

	res, err := r.readBulkString()

	if !slices.Equal(res, expected) || err != nil {
		t.Errorf("readBulkString(2\\r\\nHI\\r\\n) = %s, expected: %s, err: %v", res, expected, err)
	}
}

func TestReadBulkStringWithoutCR(t *testing.T) {
	r := initMockReader("2\n")

	res, err := r.readBulkString()

	if res != nil || err != ErrInvalidRESP {
		t.Errorf("readBulkString(2\\n) = %s, expected: %v, err: %v", res, ErrInvalidRESP, err)
	}
}

func TestReadBulkStringNonNumeric(t *testing.T) {
	r := initMockReader("asdf\r\nHI\r\n")

	res, err := r.readBulkString()

	if res != nil || err != ErrInvalidIntVal {
		t.Errorf("readBulkString(asdf\\r\\nHI\\r\\n) = %s, expected: %v, err: %v", res, ErrInvalidIntVal, err)
	}
}

func TestReadBulkStringNegativeDataLength(t *testing.T) {
	r := initMockReader("-23\r\nh\r\n")

	res, err := r.readBulkString()

	if res != nil || err != ErrIncorrectDataLen {
		t.Errorf("readBulkString(-23\\r\\nhi\\r\\n) = %s, expected: %v, err: %v", res, ErrIncorrectDataLen, err)
	}
}

func TestReadBulkStringNoNewlineAfterDataLength(t *testing.T) {
	r := initMockReader("2\rHI\r\n")

	res, err := r.readBulkString()

	if res != nil || err != ErrInvalidRESP {
		t.Errorf("readBulkString(2\\rHI\\r\\n) = %s, expected: %v, err: %v", res, ErrInvalidRESP, err)
	}
}

func TestReadBulkStringLenLessThanData(t *testing.T) {
	r := initMockReader("2\r\nHIIIII\r\n")

	res, err := r.readBulkString()

	if res != nil || err != ErrBulkStrLenMismatch {
		t.Errorf("readBulkString(2\\r\\nHIIIII\\r\\n) = %s, expected: %v, err: %v", res, ErrBulkStrLenMismatch, err)
	}
}

func TestReadBulkStringLenMoreThanData(t *testing.T) {
	r := initMockReader("455\r\nHIIIII\r\n")

	res, err := r.readBulkString()

	if res != nil || err != ErrBulkStrLenMismatch {
		t.Errorf("readBulkString(455\\r\\nHIIIII\\r\\n) = %s, expected: %v, err: %v", res, ErrBulkStrLenMismatch, err)
	}
}

func TestReadBulkStringNoEndCRLF(t *testing.T) {
	r := initMockReader("2\r\nHI")

	res, err := r.readBulkString()

	if res != nil || err != ErrInvalidRESP {
		t.Errorf("readBulkString(2\\r\\nHI) = %s, expected: %v, err: %v", res, ErrInvalidRESP, err)
	}
}

// ================================ readInt ========================================
func TestReadIntCorrectValue(t *testing.T) {
	r := initMockReader("1024\r\n")

	expected := []byte{0x0, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

	res, err := r.readInt()

	if !slices.Equal(res, expected) || err != nil {
		t.Errorf("readInt(1024\\r\\n) = %v, expected: %v, err %v", res, expected, err)
	}
}

func TestReadIntNonNumeric(t *testing.T) {
	r := initMockReader("asdf\r\n")

	res, err := r.readInt()

	if res != nil || err != ErrInvalidIntVal {
		t.Errorf("readInt(asdf\\r\\n) = %s, expected: %v, err: %v", res, ErrInvalidIntVal, err)
	}
}

func TestReadIntWithoutCR(t *testing.T) {
	r := initMockReader("2\n")

	res, err := r.readInt()

	if res != nil || err != ErrInvalidRESP {
		t.Errorf("readInt(2\\n) = %s, expected: %v, err: %v", res, ErrInvalidRESP, err)
	}
}

func TestReadIntWithoutNewLine(t *testing.T) {
	r := initMockReader("2\r")

	res, err := r.readInt()

	if res != nil || err != ErrInvalidRESP {
		t.Errorf("readInt(2\\r) = %s, expected: %v, err: %v", res, ErrInvalidRESP, err)
	}
}

// ================================ readInt ========================================
func TestReadBoolTrue(t *testing.T) {
	r := initMockReader("true\r\n")

	expected := []byte{0x01}

	res, err := r.readBool()

	if !slices.Equal(res, expected) || err != nil {
		t.Errorf("readBool(true\\r\\n) = %v, expected: %v, err %v", res, expected, err)
	}
}

func TestReadBoolFalse(t *testing.T) {
	r := initMockReader("false\r\n")

	expected := []byte{0x00}

	res, err := r.readBool()

	if !slices.Equal(res, expected) || err != nil {
		t.Errorf("readBool(false\\r\\n) = %v, expected: %v, err %v", res, expected, err)
	}
}

func TestReadBoolInvalidValue(t *testing.T) {
	r := initMockReader("asdf\r\n")

	res, err := r.readBool()

	if res != nil || err != ErrInvalidBoolVal {
		t.Errorf("readBool(asdf\\r\\n) = %s, expected: %v, err: %v", res, ErrInvalidBoolVal, err)
	}
}

func TestReadBoolWithoutCR(t *testing.T) {
	r := initMockReader("true\n")

	res, err := r.readBool()

	if res != nil || err != ErrInvalidRESP {
		t.Errorf("readBool(true\\n) = %s, expected: %v, err: %v", res, ErrInvalidRESP, err)
	}
}

func TestReadBoolWithoutNewLine(t *testing.T) {
	r := initMockReader("true\r")

	res, err := r.readBool()

	if res != nil || err != ErrInvalidRESP {
		t.Errorf("readBool(true\\r) = %s, expected: %v, err: %v", res, ErrInvalidRESP, err)
	}
}

// ================================ readArgs ========================================
func TestReadArgsCorrectArgs(t *testing.T) {
	r := initMockReader("$2\r\nHI\r\n$2\r\nYO\r\n")

	res, err := r.readArgs()
	expectedArg1 := KvsValue{dtype: BulkStrSymbol, value: []byte("HI")}
	expectedArg2 := KvsValue{dtype: BulkStrSymbol, value: []byte("YO")}

	expectedRes := []*KvsValue{&expectedArg1, &expectedArg2}

	if len(res) != len(expectedRes) || err != nil {
		t.Errorf("readArgs($2\\r\\nHI\\r\\n$2\\r\\nYO) length - %v, expected: %v, err: %v", len(res), len(expectedRes), err)
		return
	}

	for i := range expectedRes {
		expected := expectedRes[i]
		actual := res[i]

		if actual.dtype != expected.dtype || !slices.Equal(actual.value, expected.value) {
			t.Errorf("readArgs($2\\r\\nHI\\r\\n$2\\r\\nYO) value  - %v, expected: %v, err: %v", actual, expected, err)
		}
	}
}

func TestReadArgsIncorrectDtype(t *testing.T) {
	r := initMockReader("&23\r\n")

	res, err := r.readArgs()

	if res != nil || err != ErrDtypeNotSupported {
		t.Errorf("readArgs(\\r\\n) = %v, expected: %v, err: %v", res, ErrDtypeNotSupported, err)
	}
}

func TestReadArgsEmpty(t *testing.T) {
	r := initMockReader("")

	res, err := r.readArgs()

	if res != nil || err != nil {
		t.Errorf("readArgs('') = %v, expected: %v, err: %v", res, nil, err)
	}
}
