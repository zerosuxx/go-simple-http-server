package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/zerosuxx/go-simple-http-server/pkg/handler"
)

func main() {
    rootPath := "/"
    if len(os.Args) > 1 {
        rootPath = os.Args[1]
    }

    if rootPath == "-" {
        inputData, _ := io.ReadAll(os.Stdin)
        http.HandleFunc("/", handler.StdinHandler{Input: string(inputData)}.Handle)
    } else {
        http.HandleFunc("/", handler.FileHandler{RootPath: rootPath}.Handle)
    }

    log.Printf("Listening on %s | RootPath: '%s'", "http://localhost:8080", rootPath)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
