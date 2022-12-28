package base

import (
	"net/http"
	"os"
)

func ReadyHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat(readyFile); err == nil {
		return
	}
	w.WriteHeader(http.StatusNotFound)
}
