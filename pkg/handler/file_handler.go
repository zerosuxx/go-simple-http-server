package handler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type FileHandler struct {
	RootPath              string
	DirectoryIndexEnabled bool
	NotFoundFile          string
}

func (h FileHandler) handleDirectoryList(path string, w http.ResponseWriter) {
	fileList, err := getFileList(h.RootPath, path)

	if err != nil {
		handleError(err, http.StatusBadRequest, w)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`<!DOCTYPE html>` + "\n"))
	_, _ = w.Write([]byte(`<html>` + "\n"))
	_, _ = w.Write([]byte(`<head>` + "\n"))
	_, _ = w.Write([]byte(`<meta charset="utf-8">` + "\n"))
	_, _ = w.Write([]byte(`<base href="/">` + "\n"))
	_, _ = w.Write([]byte(`</head>` + "\n"))
	_, _ = w.Write([]byte(`<body>` + "\n"))
	_, _ = w.Write([]byte(fileList))
	_, _ = w.Write([]byte(`</body>` + "\n"))
	_, _ = w.Write([]byte(`</html>` + "\n"))
}

func (h FileHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var filePath = h.RootPath + "/" + r.URL.Path[1:]

	log.Printf("Incoming request path: '%s' | file: '%s'", r.URL.Path, filePath)

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if h.NotFoundFile == "" {
			handleError(err, http.StatusNotFound, w)
			return
		}

		fileInfo, err = os.Stat(h.NotFoundFile)
		if err != nil {
			log.Println("Error loading file: " + h.NotFoundFile)

			handleError(err, http.StatusNotFound, w)
			return
		}
		filePath = h.NotFoundFile
	}

	if fileInfo.IsDir() {
		if !h.DirectoryIndexEnabled {
			err := errors.New("requested file is a directory")
			handleError(err, http.StatusBadRequest, w)
			return
		}

		h.handleDirectoryList(filePath, w)

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
			_, _ = w.Write(buf[:n])
		}
	}
}

func handleError(err error, statusCode int, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	_, _ = fmt.Fprint(w, err.Error())
}

func getFileHeader(f *os.File) []byte {
	buf := make([]byte, 512)
	_, _ = f.Read(buf)
	_, _ = f.Seek(0, 0)
	return buf
}

func getFileList(rootPath string, path string) (string, error) {
	fileList := ""

	files, err := os.ReadDir(path)
	if err != nil {
		return fileList, err
	}

	urlPrefixTrimLength := len(rootPath) + 1

	for _, file := range files {
		url := path[urlPrefixTrimLength:] + "/" + file.Name()
		fileType := "file"
		if file.IsDir() {
			fileType = "dir"
		}
		fileInfo, _ := file.Info()
		fileSize := strconv.FormatInt(fileInfo.Size(), 10)

		fileList += `<a href="` + url + `">` + file.Name() + `</a> (` + fileSize + `B) [` + fileType + `]<br>` + "\n"
	}

	return fileList, err
}
