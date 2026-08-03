package middleware

import (
	"context"
	"net/http"

	"github.com/studentinovisad/popisomator/backend/internal/repository"
	"github.com/studentinovisad/popisomator/backend/internal/service"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRoles(roles ...repository.UserRole) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, ok := r.Context().Value("userID").(int64)
			if !ok {
				http.Error(w, "user ID not found in context", http.StatusInternalServerError)
				return
			}

			user, err := service.GetUserDetails(r.Context(), id)
			if err != nil {
				http.Error(w, "error fetching user details", http.StatusInternalServerError)
				return
			}

			for _, role := range roles {
				if user.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "forbidden", http.StatusForbidden)
		})
	}
}
