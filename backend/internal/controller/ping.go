package controller

import (
	"net/http"
	"io"
)

// Ping godoc
// @Summary Liveness probe
// @Tags Operational
// @Produce plain
// @Success 200 {string} string "Pong!"
// @Router /ping [get]
func Ping(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Pong!")
}