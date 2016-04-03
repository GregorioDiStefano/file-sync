package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

type FileData struct {
	FileSize int64
	FileName string
	FileHash string
}

type Files struct {
	FilePath map[string]FileData
	RootPath string
}

func ComputeMd5(filePath string, count int64) (string, error) {

	file, err := os.Open(filePath)
	hash := md5.New()
	defer file.Close()

	if err != nil {
		return "", err
	}

	if count == 0 {
		if _, err := io.Copy(hash, file); err != nil {
			return "", err
		}
	} else if _, err := io.CopyN(hash, file, count); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum([]byte{})), nil
}

func CompareDirectories(root string) Files {
	var count int64 = 1
	files := Files{FilePath: make(map[string]FileData)}
	files.RootPath = root

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if fileHash, err := ComputeMd5(path, 0); err == nil {
			files.FilePath[string(count)] = FileData{info.Size(), path[len(root):], fileHash}
			count += 1
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Failed scanning :%s, error: %s", root, err)
	}

	return files
}
