package main

import (
	"io"
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error setting up tcp listener: ", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting a new client: ", err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	buf := make([]byte, 1025) // 1025 because max size of message is 1024 bytes and we can check whether more bytes were read

	for {
		n, err := c.Read(buf)
		if err != nil {
			if err == io.EOF {
				c.Close()
				break
			}
			log.Println("Error on reading from socket: ", err)
		}

		if n == 1025 {
			c.Write([]byte(MaxCommandSizeExceeededResponse))
		}

		response := processCommand(string(buf))
		clear(buf)

		c.Write([]byte(response))
	}
}
