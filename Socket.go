package main

import (
	"bufio"
	"fmt"
	"net"
)

var conn net.Conn
var err error

const (
	PAYLOAD_PREFIX = 9
)

func sentToSocket(dst string, data []byte) bool {

	if conn == nil {
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:8081", dst))
	}

	if err != nil {
		log.Fatal("Error: ", err)
		return false
	}

	log.Debug("Socket opened.")
	fmt.Fprintf(conn, "%s", data)
	return true
}

func readFromSocket() []byte {
	connbuffer := bufio.NewReader(conn)
	var buffer []byte

	for {
		data := make([]byte, 1024)
		c, err := connbuffer.Read(data)

		//TODO: this is incorrect
		if c == 0 && len(buffer) > 0 || err != nil {
			fmt.Println("break")
			break
		} else if c > 0 {
			fmt.Println(c)
			buffer = append(buffer, data...)
			if getCompletePayload(buffer) {
				log.Infof("Recieved payload: %s", buffer[PAYLOAD_PREFIX:])
				return buffer
			}
		}

	}
	fmt.Println(buffer)
	return []byte{}
}

func getCompletePayload(buffer []byte) bool {
	if len(buffer) < PAYLOAD_PREFIX {
		return false
	}

	var meta byte
	var key uint32
	var length uint32

	log.Debugf("Payload prefix: %x\n", buffer[0:PAYLOAD_PREFIX])
	fmt.Sscanf(fmt.Sprintf("%x", buffer[0:PAYLOAD_PREFIX]),
		"%02x%08x%08x",
		&meta,
		&key,
		&length)

	if uint32(len(buffer[PAYLOAD_PREFIX:])) >= length {
		return true
	}
	return false
}
