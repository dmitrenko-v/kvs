package main

import "errors"

const ErrorSymbol = '-'

var (
	ErrServerSide           = errors.New(string(ErrorSymbol) + "ERR unexpected error on server side occured" + CRLF)
	ErrCmdNotArray          = errors.New(string(ErrorSymbol) + "ERR command must be a RESP array" + CRLF)
	ErrArrElemCountMismatch = errors.New(string(ErrorSymbol) + "ERR mismatch between number of array elements" + CRLF)
	ErrIncorrectDataLen     = errors.New(string(ErrorSymbol) + "ERR data length must be non-negative non-zero integer" + CRLF)
	ErrBulkStrLenMismatch   = errors.New(string(ErrorSymbol) + "ERR bulk string length is not correct" + CRLF)
	ErrInvalidRESP          = errors.New(string(ErrorSymbol) + "ERR invalid RESP" + CRLF)
	ErrInvalidSetCmd        = errors.New(string(ErrorSymbol) + "ERR SET command must contain key and value" + CRLF)
	ErrCmdNotSupported      = errors.New(string(ErrorSymbol) + "ERR unsupported command. Supported commands: SET, GET, PING" + CRLF)
	ErrDtypeNotSupported    = errors.New(string(ErrorSymbol) + "ERR unsupported data type. Supported types: integer, boolean, bulk string" + CRLF)
	ErrInvalidIntVal        = errors.New(string(ErrorSymbol) + "ERR Invalid integer value" + CRLF)
	ErrInvalidBoolVal       = errors.New(string(ErrorSymbol) + "ERR Invalid boolean value" + CRLF)
	ErrWrongSetArgsCount    = errors.New(string(ErrorSymbol) + "ERR SET command requires 2 args: key and value" + CRLF)
	ErrWrongGetArgsCount    = errors.New(string(ErrorSymbol) + "ERR GET command requires 1 arg: key" + CRLF)
	ErrWrongDelArgsCount    = errors.New(string(ErrorSymbol) + "ERR DELETE command requires 1 arg: key" + CRLF)
	ErrWrongKeyDtype        = errors.New(string(ErrorSymbol) + "ERR key datatype must be bulk string" + CRLF)
	ErrKeyNotExist          = errors.New(string(ErrorSymbol) + "ERR key does not exist" + CRLF)
)
