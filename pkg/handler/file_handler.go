package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type FileHandler struct {
  RootPath string
}

func (h FileHandler) Handle(w http.ResponseWriter, r *http.Request) {
  var filePath = h.RootPath + "/" + r.URL.Path[1:]

  log.Printf("Incoming request path: '%s' | file: '%s'", r.URL.Path, filePath)

	contents, err := os.ReadFile(filePath)
  if err != nil {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprint(w, err.Error())

    return
  }

  w.Header().Add("Content-Type", http.DetectContentType(contents))
  w.WriteHeader(http.StatusOK)

  w.Write(contents)
}