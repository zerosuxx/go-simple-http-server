package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/zerosuxx/go-simple-http-server/pkg/handler"
)

var Version = "development"

func main() {
    rootPath := "/"
    if len(os.Args) > 1 {
        rootPath = os.Args[1]
    }

    if rootPath == "-" {
        inputData, _ := io.ReadAll(os.Stdin)
        http.HandleFunc("/", handler.StdinHandler{Input: string(inputData)}.Handle)
    } else {
        http.HandleFunc("/", handler.FileHandler{
            RootPath: rootPath,
            DirectoryIndexEnabled: os.Getenv("DIRECTORY_INDEX_ENABLED") == "1",
        }.Handle)
    }

    log.Printf("Simple HTTP Server %s | Listening on %s | RootPath: '%s'", Version, "http://localhost:8080", rootPath)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
