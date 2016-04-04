package main

import (
	"flag"
	"math/rand"
	"testing"
)

func RandomKey() uint32 {
	return rand.Uint32()
}

/*
func TestPayload(t *testing.T) {
	br := BinaryRequest{dst: "127.0.0.1"}
	br.sendData([]byte(strings.Repeat("abc", 100)), RandomKey())

	data := ([]byte("Test"))
	br.sendData(data, RandomKey())
}

func TestSendFileWithRandom(t *testing.T) {
	isTesting = true
	br := BinaryRequest{dst: "127.0.0.1"}
	key := int32(0x00)
	fmt.Println("Random key: ", key)
	br.sendFile("tests/15_bytes", key)
}
*/
func TestSendFile(t *testing.T) {
	isTesting = true
	br := BinaryRequest{dst: "127.0.0.1"}
	br.sendFile("tests/15_bytes", RandomKey())
	br.sendFile("tests/20_bytes", RandomKey())
	br.sendFile("tests/21_bytes", RandomKey())
	br.sendFile("tests/many_0x20", RandomKey())
	//br.sendFile("tests/1M", RandomKey())
}

func TestSendFileCompressed(t *testing.T) {
	isTesting = false
	flag.Set("compress", "true")
	br := BinaryRequest{dst: "127.0.0.1"}
	br.sendFile("tests/15_bytes", RandomKey())
	br.sendFile("tests/20_bytes", RandomKey())
	br.sendFile("tests/21_bytes", RandomKey())
	br.sendFile("tests/many_0x20", RandomKey())
	br.sendFile("tests/1M", RandomKey())
}
