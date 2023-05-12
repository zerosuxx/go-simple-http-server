package handler

import (
	"net/http"
)

type StdinHandler struct {
	Input string
}

func (h StdinHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "at")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(h.Input))
}
