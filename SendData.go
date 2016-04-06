package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/bkaradzic/go-lz4"
)

type BinaryRequest struct {
	dst string
}

const (
	MAX_CHUNK        = 1024 * 1024 * 3
	COMPRESSED       = 1 << 0 //sent from server
	JSON_FILES       = 1 << 1
	JSON_LOCAL_FILES = 1 << 2 //sent from client
	FILE_PAYLOAD     = 1 << 3

	CHUNK_OVERHEAD = 9
)

func (br BinaryRequest) sendFile(filename string, key uint32) bool {
	f, err := os.Open(filename)

	if err != nil {
		log.Fatalf("Error reading file: %s", filename)
		return false
	}

	fStat, _ := f.Stat()
	fSize := fStat.Size()
	totalChunks := fSize / MAX_CHUNK
	lastPacket := fSize % MAX_CHUNK

	var i int64

	for i = 0; i < totalChunks; i++ {
		buffer := make([]byte, MAX_CHUNK)
		f.ReadAt(buffer, i*MAX_CHUNK)
		br.sendData(buffer, FILE_PAYLOAD, key)
	}

	if lastPacket > 0 {
		buffer := make([]byte, lastPacket)
		f.ReadAt(buffer, i*MAX_CHUNK)
		br.sendData(buffer, FILE_PAYLOAD, key)
	}

	return true
}

func preparePayload(meta uint8, key uint32, data []byte) []byte {
	prefix := ""

	if meta&COMPRESSED == 1 {
		fmt.Println("This chunk is compressed.")
		prefix = fmt.Sprintf("%02x%08x%08x", meta, key, len(data))
	} else {
		prefix = fmt.Sprintf("%02x%08x%08x", meta, key, len(data))
	}

	dataStr := prefix + hex.EncodeToString(data)
	dataBytes, _ := hex.DecodeString(dataStr)

	log.Debugf("Prefix: %s", []byte(prefix))
	log.Debugf("Payload: %d", len(dataBytes))

	return dataBytes
}

func isCompressible(chunk []byte) bool {
	testSize := 1000
	start := len(chunk) / 2
	end := start + testSize

	if end <= len(chunk) {
		compressed, err := lz4.Encode(nil, chunk[start:end])
		if len(compressed)+CHUNK_OVERHEAD < testSize && err == nil {
			return true
		}
	} else {
		log.Debug("Unable to check compressibility")
	}
	return false
}

func (br BinaryRequest) sendData(payload []byte, meta uint8, key uint32) []byte {
	var finalPayload []byte

	if meta == JSON_FILES || isCompressible(payload) {
		meta |= COMPRESSED
		finalPayload, _ = lz4.Encode(nil, payload)

		fmt.Printf("original: %d meta: %b, ratio: %f\n",
			len(payload),
			len(finalPayload),
			float32(len(finalPayload))/float32(len(payload))*100)

		finalPayload = preparePayload(meta, key, finalPayload)
	} else {
		finalPayload = preparePayload(meta, key, payload)
	}
	sentToSocket(br.dst, finalPayload)
	return finalPayload
}
