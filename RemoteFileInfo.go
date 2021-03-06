package main

import (
	"encoding/json"
	"fmt"
)

type RemoteFilesInfo struct {
	Status   string
	FileName string
}

type RemoteFiles interface {
	getRemoteFilesInfo(data []byte) []RemoteFilesInfo
}

func getRemoteFilesInfo(data []byte) []RemoteFilesInfo {
	var rf []RemoteFilesInfo

	if err := json.Unmarshal(data, &rf); err != nil {
		log.Critical("Unable to parse JSON data.")
		return nil
	}

	if len(rf) == 0 {
		return nil
	}

	return rf
}

func (rt RemoteFilesInfo) isCompleted(filename string) bool {

	if rt.Status == "complete" && filename == rt.FileName {
		fmt.Println("Skipping file: ", rt.FileName)
		return true
	}
	return false
}
