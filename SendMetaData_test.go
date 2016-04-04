package main

import (
	"flag"
	"os"
	"strconv"
	"testing"
)

func TestMetaDataSending(t *testing.T) {
	isTesting = false

	flag.Set("compress", "true")

	files := Files{FilePath: make(map[string]FileData)}
	/*
		filename := "tests/1M"
		fStat, _ := os.Stat(filename)
		fHash, _ := ComputeMd5(filename, 0)
		key := strconv.Itoa(1)
		files.FilePath[key] = FileData{fStat.Size(), "tests/1M", fHash}

		filename = "tests/2M"
		fStat, _ = os.Stat(filename)
		fHash, _ = ComputeMd5(filename, 0)
		key = strconv.Itoa(2)
		files.FilePath[key] = FileData{fStat.Size(), "tests/2M", fHash}
	*/
	filename := "tests/1M"
	fStat, _ := os.Stat(filename)
	fHash, _ := ComputeMd5(filename, 0)
	key := strconv.Itoa(1)
	files.FilePath[key] = FileData{fStat.Size(), "tests/1M", fHash}

	filename = "tests/2M"
	fStat, _ = os.Stat(filename)
	fHash, _ = ComputeMd5(filename, 0)
	key = strconv.Itoa(2)
	files.FilePath[key] = FileData{fStat.Size(), "tests/2M", fHash}

	filename = "tests/3M"
	fStat, _ = os.Stat(filename)
	fHash, _ = ComputeMd5(filename, 0)
	key = strconv.Itoa(3)
	files.FilePath[key] = FileData{fStat.Size(), "tests/3M", fHash}

	filename = "tests/2GB"
	fStat, _ = os.Stat(filename)
	fHash, _ = ComputeMd5(filename, 0)
	key = strconv.Itoa(4)
	files.FilePath[key] = FileData{fStat.Size(), "tests/2GB", fHash}
	/*
			filename = "tests/3M"
			fStat, _ = os.Stat(filename)
			fHash, _ = ComputeMd5(filename, 0)
			key = strconv.Itoa(3)
			files.FilePath[key] = FileData{fStat.Size(), "tests/3M", fHash}
			/*
		filename = "tests/lots_of_numbers"
		fStat, _ = os.Stat(filename)
		fHash, _ = ComputeMd5(filename, 0)
		key = strconv.Itoa(4)
		files.FilePath[key] = FileData{fStat.Size(), "tests/lots_of_numbers", fHash}
	*/
	br := BinaryRequest{dst: "127.0.0.1"}
	br.sendFileInfo(files)
	readFromSocket()
	br.sendFile("tests/1M", uint32(1))
	br.sendFile("tests/2M", uint32(2))
	br.sendFile("tests/3M", uint32(3))
	br.sendFile("tests/2GB", uint32(4))
	//br.sendFile("tests/lots_of_numbers", int32(4))
}
