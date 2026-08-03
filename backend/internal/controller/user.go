package controller

import (
	"encoding/json"
	"net/http"

	"github.com/studentinovisad/popisomator/backend/internal/service"
)

func UserDetailsPersonal(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "not logged in", http.StatusForbidden)
		return
	}

	id, err := service.ValidateToken(cookie.Value)
	if err != nil {
		http.Error(w, "invalid token", http.StatusForbidden)
		return
	}

	user, err := service.GetUserDetails(r.Context(), id)
	if err != nil {
		http.Error(w, "error fetching details", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
