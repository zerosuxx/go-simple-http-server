package handler

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
)

type FileHandler struct {
	RootPath string
}

type Emitter struct {
	Writer http.ResponseWriter
  FileInfo fs.FileInfo
	header string
}

func (e *Emitter) Emit(data []byte) {
	if e.header == "" {
		
	}

	e.Writer.Write(data)
}

func handleError(err error, statusCode int, w http.ResponseWriter) {
  w.WriteHeader(http.StatusNotFound)
  fmt.Fprint(w, err.Error())
}

func getFileHeader(f *os.File) []byte {
  buf := make([]byte, 512)
	_, _ = f.Read(buf)
  f.Seek(0, 0)
  return buf
}

func (h FileHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var filePath = h.RootPath + "/" + r.URL.Path[1:]

	log.Printf("Incoming request path: '%s' | file: '%s'", r.URL.Path, filePath)

  fileInfo, err := os.Stat(filePath)
  if err != nil {
    handleError(err, http.StatusNotFound, w)

    return
  }
  if fileInfo.IsDir() {
    err := errors.New("requested file is a directory")
    handleError(err, http.StatusBadRequest, w)

    return
  }
  
	f, err := os.Open(filePath)
	if err != nil {
    handleError(err, http.StatusInternalServerError, w)

		return
	}
	defer f.Close()

  w.Header().Add("Content-Type", http.DetectContentType(getFileHeader(f)))
  w.Header().Add("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
  w.WriteHeader(http.StatusOK)
  
	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		if n > 0 {
      w.Write(buf[:n])
		}
	}
}
