package main

import (
	"bufio"
	"encoding/binary"
	"io"
	"slices"
	"strconv"
)

const (
	ArrSymbol       = '*'
	SimpleStrSymbol = '+'
	BulkStrSymbol   = '$'
	IntSymbol       = ':'
	BoolSymbol      = '#'
)

type respReader struct {
	reader *bufio.Reader
}

func NewRespReader(r io.Reader) *respReader {
	return &respReader{reader: bufio.NewReader(r)}
}

// EVERY RESP command MUST start like "*<no. of lines after first one>\r\n$" and then input may vary
// so this function basically checks this beggining and then reads args
func (r *respReader) readCommand() (command string, args []*KvsValue, err error) {
	firstLine, err := r.readRespLine()
	if err != nil {
		return "", nil, err
	}

	if err := readArray(firstLine); err != nil {
		return "", nil, err
	}

	curByte, err := r.reader.ReadByte()
	if err != nil || curByte != BulkStrSymbol {
		return "", nil, ErrInvalidRESP
	}

	cmd, err := r.readBulkString()
	if err != nil {
		return "", nil, err
	}

	args, err = r.readArgs()
	if err != nil {
		return "", nil, err
	}

	return string(cmd), args, nil
}

func (r *respReader) readArgs() (args []*KvsValue, err error) {
	_, err = r.reader.Peek(1) // checking for empty args
	if err != nil {
		return nil, nil
	}

	maxArgsCount := 2
	args = make([]*KvsValue, 0, maxArgsCount)
	curArgIx := 0

	for {
		if curArgIx == maxArgsCount {
			break
		}

		dtype, err := r.reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, ErrInvalidRESP
		}

		var val []byte
		switch dtype {
		case BulkStrSymbol:
			val, err = r.readBulkString()
		case BoolSymbol:
			val, err = r.readBool()
		case IntSymbol:
			val, err = r.readInt()
		default:
			return nil, ErrDtypeNotSupported
		}

		if err != nil {
			return nil, err
		}

		arg := &KvsValue{dtype: dtype, value: val}
		args = append(args, arg)

		curArgIx++
	}

	return args, nil
}

func (r *respReader) readRespLine() (line []byte, err error) {
	respLine, err := r.reader.ReadBytes('\n')

	if err != nil {
		if err == io.EOF {
			return nil, io.EOF
		}

		return nil, ErrInvalidRESP
	}

	crIndex := len(respLine) - 2
	if respLine[crIndex] != '\r' {
		return nil, ErrInvalidRESP
	}

	return respLine[:crIndex], nil
}

// For this minimum version of kvs we dont really need arrays except just read command header so it differs from other read<type> funcs
func readArray(respFirstLine []byte) error {
	if respFirstLine[0] != ArrSymbol {
		return ErrCmdNotArray
	}

	dataLength, err := bytesToInt(respFirstLine[1:])
	if err != nil {
		return err
	}

	if dataLength <= 0 {
		return ErrIncorrectDataLen
	}

	return nil
}

func (r *respReader) readBulkString() (val []byte, err error) {
	dataLengthBytes, err := r.reader.ReadBytes('\r')
	if err != nil {
		return nil, ErrInvalidRESP
	}

	dataLength, err := bytesToInt(dataLengthBytes[:len(dataLengthBytes)-1])
	if err != nil {
		return nil, err
	}

	if dataLength <= 0 {
		return nil, ErrIncorrectDataLen
	}

	lastByte, err := r.reader.ReadByte()
	if err != nil || lastByte != '\n' {
		return nil, ErrInvalidRESP
	}

	res := make([]byte, dataLength)

	_, err = io.ReadFull(r.reader, res)
	if err != nil {
		return nil, ErrBulkStrLenMismatch
	}

	crlf, err := r.reader.ReadBytes('\n')
	if err != nil {
		return nil, ErrInvalidRESP
	}

	if len(crlf) > 2 {
		return nil, ErrBulkStrLenMismatch
	}

	if !slices.Equal(crlf, []byte(CRLF)) {
		return nil, ErrInvalidRESP
	}

	return res, nil
}

func (r *respReader) readBool() (val []byte, err error) {
	valBytes, err := r.reader.ReadBytes('\r')
	if err != nil {
		return nil, ErrInvalidRESP
	}

	valBool, err := strconv.ParseBool(string(valBytes[:len(valBytes)-1]))
	if err != nil {
		return nil, ErrInvalidBoolVal
	}

	lastByte, err := r.reader.ReadByte()
	if err != nil || lastByte != '\n' {
		return nil, ErrInvalidRESP
	}

	res := make([]byte, 1)

	if valBool {
		res[0] = 0x01
	} else {
		res[0] = 0x00
	}

	return res, nil
}

// I thought that fucking redis-cli  send int data as int datatype, BUT ITS STRING. WHY????
func (r *respReader) readInt() (val []byte, err error) {
	valBytes, err := r.reader.ReadBytes('\r')
	if err != nil {
		return nil, ErrInvalidRESP
	}

	bytesLen := len(valBytes)

	intVal, err := bytesToInt(valBytes[:bytesLen-1])
	if err != nil {
		return nil, err
	}

	lastByte, err := r.reader.ReadByte()
	if err != nil || lastByte != '\n' {
		return nil, ErrInvalidRESP
	}

	res := make([]byte, 8)
	binary.NativeEndian.PutUint64(res, uint64(intVal))

	return res, nil
}
