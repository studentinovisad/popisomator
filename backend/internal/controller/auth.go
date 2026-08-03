package controller

import (
	"encoding/json"
	"net/http"

	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/service"
)

func Login(w http.ResponseWriter, r *http.Request) {
	body := http.MaxBytesReader(w, r.Body, 1024)

	var req dto.LoginRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	token, err := service.Login(r.Context(), req)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusBadRequest)
		return
	}

	cookie := http.Cookie{
		Name:     "session",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:   "session",
		MaxAge: -1,
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}
