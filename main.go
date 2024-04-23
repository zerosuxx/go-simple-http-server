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
	rootPath := "."
	if len(os.Args) > 1 {
		rootPath = os.Args[1]
	}

	if rootPath == "-" {
		inputData, _ := io.ReadAll(os.Stdin)
		http.HandleFunc("/", handler.StdinHandler{Input: string(inputData)}.Handle)
	} else {
		http.HandleFunc("/", handler.FileHandler{
			RootPath:              rootPath,
			DirectoryIndexEnabled: os.Getenv("DIRECTORY_INDEX_ENABLED") == "1",
			NotFoundFile:          os.Getenv("NOT_FOUND_FILE"),
		}.Handle)
	}

	address := os.Getenv("ADDRESS")
	if address == "" {
		address = "0.0.0.0:8080"
	}
	log.Printf("Simple HTTP Server %s | Listening on %s | RootPath: '%s'", Version, "http://"+address, rootPath)
	log.Fatal(http.ListenAndServe(address, nil))
}
