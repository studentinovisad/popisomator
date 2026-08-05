package middleware

import (
	"context"
	"net/http"

	"github.com/studentinovisad/popisomator/backend/internal/repository"
	"github.com/studentinovisad/popisomator/backend/internal/response"
	"github.com/studentinovisad/popisomator/backend/internal/service"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			response.WriteError(w, http.StatusUnauthorized, "not logged in")
			return
		}

		id, err := service.ValidateToken(cookie.Value)
		if err != nil {
			response.WriteError(w, http.StatusUnauthorized, "invalid token")
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
				response.WriteError(w, http.StatusInternalServerError, "user ID not found in context")
				return
			}

			user, err := service.GetUserDetails(r.Context(), id)
			if err != nil {
				response.WriteError(w, http.StatusInternalServerError, "error fetching user details")
				return
			}

			for _, role := range roles {
				if user.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			response.WriteError(w, http.StatusForbidden, "forbidden")
		})
	}
}
