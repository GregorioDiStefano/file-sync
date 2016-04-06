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
	backend1Leveled.SetLevel(logging.DEBUG, "")
	logging.SetBackend(backend1Leveled)

	flag.String("host", "127.0.0.1", "host to connect to")
	flag.String("dir", "", "directory to send files from")
	flag.Parse()
}

func main() {
	local := CompareDirectories(flag.Lookup("dir").Value.String())
	br := BinaryRequest{dst: flag.Lookup("host").Value.String()}
	br.sendFileInfo(local)
	readFromSocket()

	for fileKey, fileValue := range local.FilePath {
		fmt.Println("Sending: ", local.RootPath+fileValue.FileName)
		fileKey, _ := strconv.ParseUint(fileKey, 10, 32)
		br.sendFile(local.RootPath+fileValue.FileName, uint32(fileKey))
	}
}
