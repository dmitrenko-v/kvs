package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

const (
	SetCmd     = "SET"
	GetCmd     = "GET"
	DeleteCmd  = "DELETE"
	CommandCmd = "COMMAND" // Command exists just for correct connection through redis-cli
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Port to run application on")
	flag.Parse()

	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal("Error setting up tcp listener: ", err)
	}
	log.Println("Application started at port:", port)

	initStorage()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting a new client: ", err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()
	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic occured: ", r)
			c.Write([]byte(ErrServerSide.Error()))
		}
	}()

	respReader := NewRespReader(c)

	for {
		cmd, args, err := respReader.readCommand()
		if err != nil {
			if err == io.EOF {
				log.Println("Client disconnected")
				c.Close()
				return
			}

			c.Write([]byte(err.Error()))
			continue
		}

		switch strings.ToUpper(cmd) {
		case CommandCmd:
			c.Write([]byte(OkResponse))
		case SetCmd:
			err = setHandler(args)
			if err == nil {
				c.Write([]byte(OkResponse))
			}
		case GetCmd:
			res, err := getHandler(args)
			if res == nil && err == nil {
				c.Write([]byte(NullResponse))
			} else if res != nil && err == nil {
				response := kvsValueToResponse(res)
				c.Write(response)
			}
		case DeleteCmd:
			err = deleteHandler(args)
			if err == nil {
				c.Write([]byte(OkResponse))
			}
		default:
			c.Write([]byte(ErrCmdNotSupported.Error()))
		}

		if err != nil {
			c.Write([]byte(err.Error()))
		}
	}
}

func kvsValueToResponse(kvsValue *KvsValue) []byte {
	var sb strings.Builder
	var dataValue string

	sb.WriteByte(kvsValue.dtype)

	switch kvsValue.dtype {
	case IntSymbol:
		dataValue = decodeInt(kvsValue.value)
	case BoolSymbol:
		dataValue = decodeBool(kvsValue.value)
	default:
		dataValue = string(kvsValue.value)
		blkStrLen := len(dataValue)
		sb.WriteString(strconv.Itoa(blkStrLen))
		sb.WriteString(CRLF)
	}

	sb.WriteString(dataValue)
	sb.WriteString(CRLF)

	return []byte(sb.String())
}
