package router

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/studentinovisad/popisomator/backend/internal/config"
	"github.com/studentinovisad/popisomator/backend/internal/controller"
	_ "github.com/studentinovisad/popisomator/backend/internal/docs"
	"github.com/studentinovisad/popisomator/backend/internal/middleware"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()

	// Operational
	mux.HandleFunc("/ping", controller.Ping)
	mux.HandleFunc("/health", controller.Healthcheck)

	// Swagger UI, dev-only
	if config.CurrentConfig.DebugMode {
		mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	}

	// Auth & registration
	mux.HandleFunc("POST /auth/login", controller.Login)
	mux.HandleFunc("POST /auth/logout", controller.Logout)
	mux.HandleFunc("POST /auth/register", controller.Register)

	// Users
	mux.Handle("POST /users", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("admin"),
		middleware.Handle(controller.CreateUser),
	))
	mux.Handle("GET /users", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("admin"),
		middleware.Handle(controller.ListUsers),
	))
	mux.Handle("GET /users/me", middleware.Chain(
		middleware.RequireAuth,
		middleware.Handle(controller.UserDetailsPersonal),
	))
	mux.Handle("PATCH /users/{id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("admin"),
		middleware.Handle(controller.UpdateUser),
	))
	mux.Handle("DELETE /users/{id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("admin"),
		middleware.Handle(controller.DeleteUser),
	))

	// Items
	mux.Handle("POST /items", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.CreateItem),
	))
	mux.Handle("GET /items", middleware.Chain(
		middleware.RequireAuth,
		middleware.Handle(controller.ListItems),
	))
	mux.Handle("GET /items/{id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.Handle(controller.GetItem),
	))
	mux.Handle("PATCH /items/{id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.UpdateItem),
	))
	mux.Handle("POST /items/{id}/consume", middleware.Chain(
		middleware.RequireAuth,
		middleware.Handle(controller.ConsumeItem),
	))
	mux.Handle("DELETE /items/{id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.DeleteItem),
	))
	mux.Handle("POST /items/{id}/properties", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.AddItemProperty),
	))
	mux.Handle("PUT /items/{id}/properties/{prop_id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.UpdateItemProperty),
	))
	mux.Handle("DELETE /items/{id}/properties/{prop_id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.RemoveItemProperty),
	))

	// Properties
	mux.Handle("POST /properties", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.CreateProperty),
	))
	mux.Handle("GET /properties", middleware.Chain(
		middleware.RequireAuth,
		middleware.Handle(controller.ListProperties),
	))
	mux.Handle("GET /properties/options", middleware.Chain(
		middleware.RequireAuth,
		middleware.Handle(controller.GetPropertyOptions),
	))
	mux.Handle("GET /properties/{id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.Handle(controller.GetProperty),
	))
	mux.Handle("PATCH /properties/{id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.UpdateProperty),
	))
	mux.Handle("DELETE /properties/{id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.DeleteProperty),
	))

	// Item Types
	mux.Handle("POST /item-types", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.CreateItemType),
	))
	mux.Handle("GET /item-types", middleware.Chain(
		middleware.RequireAuth,
		middleware.Handle(controller.ListItemTypes),
	))
	mux.Handle("GET /item-types/options", middleware.Chain(
		middleware.RequireAuth,
		middleware.Handle(controller.GetItemTypeOptions),
	))
	mux.Handle("GET /item-types/{id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.Handle(controller.GetItemType),
	))
	mux.Handle("PATCH /item-types/{id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.UpdateItemType),
	))
	mux.Handle("DELETE /item-types/{id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.DeleteItemType),
	))
	mux.Handle("POST /item-types/{id}/properties", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.AddItemTypeProperty),
	))
	mux.Handle("PATCH /item-types/{id}/properties/{prop_id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.UpdateItemTypeProperty),
	))
	mux.Handle("DELETE /item-types/{id}/properties/{prop_id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.RemoveItemTypeProperty),
	))

	return mux
}
