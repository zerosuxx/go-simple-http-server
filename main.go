package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func main() {
    http.HandleFunc("/", FileHandler)

    log.Println("Listening on http://localhost:8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func FileHandler(w http.ResponseWriter, r *http.Request) {
    contents, err := os.ReadFile("/tmp/" + r.URL.Path[1:])
    if err != nil {
      w.WriteHeader(http.StatusNotFound)
      fmt.Fprintf(w, err.Error())
      return
    }

  
    mimeType := http.DetectContentType(contents)
    if mimeType == "application/octet-stream" {
       mimeType = "application/text"
    }
    w.Header().Add("Content-Type", mimeType)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(string(contents)))
}
