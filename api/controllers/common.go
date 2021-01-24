package controllers

import (
	"io"
	"net/http"
)

func sendMessageIfError(err error, w http.ResponseWriter, statusCode int, message string) (ok bool) {
	if err != nil {
		w.WriteHeader(statusCode)
		io.WriteString(w, message)
	} else {
		ok = true
	}

	return
}
