package main

import (
	"fmt"
	"github.com/op/go-logging"
	"os"
)

var isTesting bool = false
var log = logging.MustGetLogger("filesync")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func init() {
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")
	logging.SetBackend(backend1Leveled)
}

func main() {
	local := CompareDirectories("/home/greg/Desktop/file-sync/tests/test1/A")
	remote := CompareDirectories("/home/greg/Desktop/file-sync/tests/test1/B")

	for fileName, fileValue := range local.FilePath {
		if fileValue != remote.FilePath[fileName] {
			fmt.Println("Invalid remote file: ", fileName, fileValue, remote.FilePath[fileName].FileSize)
			NewFile(local.RootPath + fileName).isResumable(remote.FilePath[fileName])
		}
	}
}
