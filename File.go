package main

import (
	"fmt"
	"os"
)

type File struct {
	filePath string
	fileSize int64
	hash     string
}

//This is a struct containing information of the local file
func NewFile(filepath string) *File {
	f := new(File)
	f.filePath = filepath

	if fileInfo, err := os.Stat(filepath); err != nil {
		log.Fatal("Unable to stat: ", filepath)
	} else {
		f.fileSize = fileInfo.Size()
	}

	fmt.Printf("Hashing: %s\n", filepath)
	f.hash, _ = ComputeMd5(filepath, 0)
	return f
}

func (localFile File) isResumable(remoteFile FileData) bool {
	if localFile.fileSize < remoteFile.FileSize {
		log.Debug("Remote file is larger than local file")
		return false
	} else if remoteFile.FileSize < localFile.fileSize {
		h, _ := ComputeMd5(localFile.filePath, remoteFile.FileSize)
		log.Debug("File download can be resumable: ", h == remoteFile.FileHash)
		return h == remoteFile.FileHash
	} else {
		return false
	}
}
