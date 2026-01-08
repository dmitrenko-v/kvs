package main

const (
	MaxCommandSizeExceeededResponse = string(ErrorSymbol) + "ERR max command size is 1024 bytes" + CRLF
	ServerSideErrorResponse         = string(ErrorSymbol) + "ERR unexpected error on server side occured" + CRLF
	CommandMustBeAnArrayResponse    = string(ErrorSymbol) + "ERR command must be a RESP array" + CRLF
	ArrayElemCountMismatch          = string(ErrorSymbol) + "ERR mismatch between number of array elements" + CRLF
	MissingDollarResponse           = string(ErrorSymbol) + "ERR there must be dollar sign before data length" + CRLF
	DataLengthNonNumericResponse    = string(ErrorSymbol) + "ERR data length must be integer" + CRLF
	InvalidRESPResponse             = string(ErrorSymbol) + "ERR invalid RESP" + CRLF
	InvalidSetCommandResponse       = string(ErrorSymbol) + "ERR SET command must contain key and value" + CRLF
	UnsupportedCommandResponse      = string(ErrorSymbol) + "ERR unsupported command. Supported commands: SET, GET, PING" + CRLF
	OkResponse                      = string(SimpleStringSymbol) + "OK" + CRLF
	PongResponse                    = string(SimpleStringSymbol) + "PONG" + CRLF
)
