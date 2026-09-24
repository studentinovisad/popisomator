package middleware

import (
	"context"
	"errors"
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

func HasRoles(ctx context.Context, roles ...repository.UserRole) (bool, error) {
	id, ok := ctx.Value("userID").(int64)
	if !ok {
		return false, errors.New("user ID not found in context")
	}

	user, err := service.GetUserDetails(ctx, id)
	if err != nil {
		return false, errors.New("error fetching user details")
	}

	for _, role := range roles {
		if user.Role == role {
			return true, nil
		}
	}
	return false, nil
}

func RequireRoles(roles ...repository.UserRole) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hasRole, err := HasRoles(r.Context(), roles...)
			if err != nil {
				response.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			if hasRole {
				next.ServeHTTP(w, r)
				return
			}

			response.WriteError(w, http.StatusForbidden, "forbidden")
		})
	}
}
