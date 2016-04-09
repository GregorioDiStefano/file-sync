package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("filesync")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func init() {
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend1Leveled := logging.AddModuleLevel(backend1)
	backendFormatter := logging.NewBackendFormatter(backend1, format)
	logging.SetBackend(backend1Leveled, backendFormatter)

	flag.String("host", "127.0.0.1", "host to connect to")
	flag.String("dir", "", "directory to send files from")
	flag.Bool("debug", false, "show debug messages")
	flag.Parse()
}

func main() {
	debugSet, _ := strconv.ParseBool(flag.Lookup("debug").Value.String())
	if debugSet == true {
		logging.SetLevel(logging.DEBUG, "*")
	}

	local := CompareDirectories(flag.Lookup("dir").Value.String())
	br := BinaryRequest{dst: flag.Lookup("host").Value.String()}
	br.sendFileInfo(local)

	data := readFromSocket()
	rt := getRemoteFilesInfo(data)

	for fileKey, fileValue := range local.FilePath {
		var ignoreFile bool

		for _, data := range rt {

			if data.isCompleted(fileValue.FileName) {
				ignoreFile = true
			}

		}

		if ignoreFile {
			continue
		}

		fmt.Println("Sending: ", fileValue.FileName)
		fileKey, _ := strconv.ParseUint(fileKey, 10, 32)
		br.sendFile(fileValue.FileName, uint32(fileKey))
	}
}
