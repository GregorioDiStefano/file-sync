package main

import (
	"encoding/hex"
	"fmt"
	"github.com/bkaradzic/go-lz4"
	"os"
)

type BinaryRequest struct {
	dst string
}

const (
	MAX_CHUNK = 1024 * 1024 * 3
	JSON      = 0
)

func (br BinaryRequest) sendFile(filename string, key int32) bool {
	f, err := os.Open(filename)

	if err != nil {
		fmt.Errorf("Error reading file: ", filename)
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
		br.sendData(buffer, key)
	}

	if lastPacket > 0 {
		buffer := make([]byte, lastPacket)
		f.ReadAt(buffer, i*MAX_CHUNK)
		br.sendData(buffer, key)
	}

	return true
}

func preparePayload(compressed bool, key int32, data []byte) []byte {
	prefix := ""

	if compressed {
		fmt.Println("This chunk is compressed.")
		prefix = fmt.Sprintf("%02x%08x%08x", 1, key, len(data))
	} else {
		prefix = fmt.Sprintf("%02x%08x%08x", 0, key, len(data))
	}

	dataStr := prefix + hex.EncodeToString(data)
	dataBytes, _ := hex.DecodeString(dataStr)

	log.Debug("Prefix: %s", []byte(prefix))
	log.Debug("Payload: %d", len(dataBytes))

	return dataBytes
}

func isCompressible(chunk []byte) bool {
	testSize := 1000
	start := len(chunk) / 2
	end := start + testSize

	if end <= len(chunk) {
		compressed, err := lz4.Encode(nil, chunk[start:end])
		if len(compressed) < testSize && err == nil {
			return true
		}
	} else {
		log.Debug("Unable to check compressibility")
	}
	return false
}

func (br BinaryRequest) sendData(payload []byte, key int32) []byte {
	var finalPayload []byte

	if key == JSON || isCompressible(payload) {
		finalPayload, _ = lz4.Encode(nil, payload)

		fmt.Printf("original: %d compressed: %d, ratio: %f\n",
			len(payload),
			len(finalPayload),
			float32(len(finalPayload))/float32(len(payload))*100)

		finalPayload = preparePayload(true, key, finalPayload)
	} else {
		finalPayload = preparePayload(false, key, payload)
	}
	sentToSocket(br.dst, finalPayload)
	return finalPayload
}
