package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	SetCmd     = "SET"
	GetCmd     = "GET"
	CommandCmd = "COMMAND"
	PingCmd    = "PING"
)

// the main idea of command processing is splitting command by CRLF and have command parts as slice of strings
// then process this command parts one by one or by two if its bulk string for example
func processCommand(command string) string {
	commandParts := strings.Split(command, CRLF)
	commandParts = commandParts[:len(commandParts)-1] // its stupid but the last element for some reason is some unknown shit so removing it
	commandPartsCount := len(commandParts)

	if commandPartsCount < 3 { // command must at least consists of 3 parts: 1st part-array datatype, and bulk string for command(2 parts-length and data)
		return InvalidRESPResponse
	}

	if err := checkArrayValidity(commandPartsCount, commandParts[0]); err != nil {
		return err.Error()
	}

	cmdLen := commandParts[1]
	cmd := commandParts[2]
	if err := checkBulkStringValidity(cmdLen, cmd); err != nil {
		return err.Error()
	}

	switch cmd {
	case PingCmd:
		return PongResponse
	case CommandCmd:
		return OkResponse // just so redis-cli can be used because it sends COMMAND DOCS initially when connected
	case SetCmd:
		// key, value, err := parseSetCmd(commandParts[2:])
		// if err != nil {
		// 	return err.Error()
		// }
		return OkResponse
	default:
		return UnsupportedCommandResponse
	}
}

func parseSetCommand(setCmdParts []string) (key string, value string, err error) {
	if len(setCmdParts) < 4 { // 4 is minimum count of parts required for set command(including an actual "SET") if we split command by CRLF
		return "", "", errors.New(InvalidSetCommandResponse)
	}

	keyLen := setCmdParts[1]
	keyData := setCmdParts[2]
	if err := checkBulkStringValidity(keyLen, keyData); err != nil {
		return "", "", err
	}

	valueType := setCmdParts[3][0] // first char of the part that's following key is a value type

	switch valueType { // at this point, we just checking validity of input, actual encoding happening in storage.go
	case IntegerSymbol:
		err = checkIntValueValidity(setCmdParts[3][1:])
	case BooleanSymbol:
		err = checkBoolValueValidity(setCmdParts[3][1:])
	case BulkStringSymbol:

	}

	if err != nil {
		return "", "", err
	}

	return keyData, "", nil
}

func checkIntValueValidity(intStr string) error {
	if _, err := strconv.Atoi(intStr); err != nil {
		response := fmt.Sprintf("%vERR %v is not a valid integer", ErrorSymbol, intStr)
		return errors.New(response)
	}

	return nil
}

func checkBoolValueValidity(boolStr string) error {
	if boolStr != "true" && boolStr != "false" {
		response := fmt.Sprintf("%vERR %v is not a valid boolean", ErrorSymbol, boolStr)
		return errors.New(response)
	}

	return nil
}

func checkArrayValidity(cmdPartsCount int, dataTypeCmdPart string) error {
	if dataTypeCmdPart[0] != ArraySymbol {
		return errors.New(CommandMustBeAnArrayResponse)
	}

	elementsCount, err := strconv.Atoi(dataTypeCmdPart[1:])
	if err != nil {
		return errors.New(DataLengthNonNumericResponse)
	}

	expectedCmdPartsCount := elementsCount * 2 // *2 because each element also has line with length before it by RESP protocol

	if cmdPartsCount-1 != expectedCmdPartsCount { // -1 for array type part
		return errors.New(ArrayElemCountMismatch)
	}

	return nil
}

func checkBulkStringValidity(lengthCmdPart string, data string) error {
	if lengthCmdPart[0] != BulkStringSymbol {
		return errors.New(MissingDollarResponse)
	}

	length, err := strconv.Atoi(lengthCmdPart[1:])
	if err != nil {
		return errors.New(DataLengthNonNumericResponse)
	}

	dataLen := len(data)
	if length != dataLen {
		response := fmt.Sprintf("%vERR bulk string data and length mismatch: expected data length - %v, got - %v", ErrorSymbol, length, dataLen)
		return errors.New(response)
	}

	return nil
}
