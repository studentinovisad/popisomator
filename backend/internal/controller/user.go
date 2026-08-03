package controller

import (
	"encoding/json"
	"net/http"

	"github.com/studentinovisad/popisomator/backend/internal/service"
)

func UserDetailsPersonal(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("userID").(int64)
	if !ok {
		http.Error(w, "user ID not found in context", http.StatusInternalServerError)
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
