package main

import (
	"fmt"
	"net"
)

var conn net.Conn
var err error

func sentToSocket(dst string, data []byte) bool {

	if conn == nil {
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:8081", dst))
	}

	if err != nil {
		log.Fatal("Error: ", err)
		return false
	} else {
		log.Debug("Socket opened.")
	}

	fmt.Fprintf(conn, "%s", data)
	return true
}

func readFromSocket() int {
	var buf = make([]byte, 1024)

	for {
		readlen, err := conn.Read(buf)
		fmt.Println("Reading: ", readlen)
		if err != nil {
			log.Fatal("Error when reading from socket: %s\n", err)
			return 0
		}
		if readlen == 0 {
			log.Debug("Connection closed by remote host\n")
			return 0
		}
		fmt.Printf("Client at %s says '%s'\n", conn.RemoteAddr().String(), buf)
	}
}
