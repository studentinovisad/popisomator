package router

import (
	"net/http"

	"github.com/studentinovisad/popisomator/backend/internal/controller"
	"github.com/studentinovisad/popisomator/backend/internal/middleware"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", controller.Ping)
	mux.HandleFunc("/health", controller.Healthcheck)
	mux.HandleFunc("POST /auth/login", controller.Login)
	mux.HandleFunc("POST /auth/logout", controller.Logout)
	mux.Handle("POST /auth/register", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("admin"),
		middleware.Handle(controller.Register),
	))
	mux.Handle("GET /user/details", middleware.Chain(
		middleware.RequireAuth,
		middleware.Handle(controller.UserDetailsPersonal),
	))

	return mux
}
