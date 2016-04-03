package main

import (
	"encoding/json"
)

func (br BinaryRequest) sendFileInfo(filenameMapping Files) bool {
	data, err := json.Marshal(filenameMapping)

	if err != nil {
		log.Fatal("Error sending file meta data: ", err)
		return false
	}
	log.Debug("Sending FileInfo JSON: ", string(data))
	br.sendData(data, JSON)
	return true
}
