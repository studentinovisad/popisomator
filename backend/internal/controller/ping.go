package controller

import (
	"net/http"
	"io"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Pong!")
}