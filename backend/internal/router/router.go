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
	mux.HandleFunc("POST /auth/register", controller.Register)
	mux.Handle("POST /auth/create", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("admin"),
		middleware.Handle(controller.CreateUser),
	))
	mux.Handle("GET /user/details", middleware.Chain(
		middleware.RequireAuth,
		middleware.Handle(controller.UserDetailsPersonal),
	))
	mux.Handle("GET /users", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("admin"),
		middleware.Handle(controller.ListUsers),
	))
	mux.Handle("PATCH /user/{id}/role", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("admin"),
		middleware.Handle(controller.UpdateRole),
	))
	mux.Handle("POST /user/{id}/approve", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("admin"),
		middleware.Handle(controller.ApproveRegistration),
	))
	mux.Handle("POST /user/{id}/decline", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("admin"),
		middleware.Handle(controller.DeclineRegistration),
	))
	mux.Handle("GET /item", middleware.Chain(
		middleware.RequireAuth,
		middleware.Handle(controller.GetAllItems),
	))
	mux.Handle("GET /item/{id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.Handle(controller.GetItem),
	))
	mux.Handle("POST /item", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.CreateItem),
	))
	mux.Handle("DELETE /item/{id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.DeleteItem),
	))
	mux.Handle("POST /item/{id}/consume", middleware.Chain(
		middleware.RequireAuth,
		middleware.Handle(controller.ConsumeItem),
	))
	mux.Handle("PATCH /item/{id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.SetItemType),
	))
	mux.Handle("GET /item/properties", middleware.Chain(
		middleware.RequireAuth,
		middleware.Handle(controller.GetAllProperties),
	))
	mux.Handle("GET /item/properties/page", middleware.Chain(
		middleware.RequireAuth,
		middleware.Handle(controller.ListProperties),
	))
	mux.Handle("GET /item/properties/{id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.Handle(controller.GetProperty),
	))
	mux.Handle("POST /item/properties", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.CreateProperty),
	))
	mux.Handle("PATCH /item/properties/{id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.UpdateProperty),
	))
	mux.Handle("DELETE /item/properties/{id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.DeleteProperty),
	))
	mux.Handle("POST /item/{id}/properties", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.AddItemProperty),
	))
	mux.Handle("PUT /item/{id}/properties/{prop_id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.UpdateItemProperty),
	))
	mux.Handle("DELETE /item/{id}/properties/{prop_id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.RemoveItemProperty),
	))
	mux.Handle("GET /item/types", middleware.Chain(
		middleware.RequireAuth,
		middleware.Handle(controller.GetAllItemTypes),
	))
	mux.Handle("GET /item/types/page", middleware.Chain(
		middleware.RequireAuth,
		middleware.Handle(controller.ListItemTypes),
	))
	mux.Handle("GET /item/types/{id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.Handle(controller.GetItemType),
	))
	mux.Handle("POST /item/types", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.CreateItemType),
	))
	mux.Handle("PATCH /item/types/{id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.UpdateItemType),
	))
	mux.Handle("DELETE /item/types/{id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.DeleteItemType),
	))
	mux.Handle("POST /item/types/{id}/properties", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.AddItemTypeProperty),
	))
	mux.Handle("PATCH /item/types/{id}/properties/{prop_id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.UpdateItemTypeProperty),
	))
	mux.Handle("DELETE /item/types/{id}/properties/{prop_id}", middleware.Chain(
		middleware.RequireAuth,
		middleware.RequireRoles("manager", "admin"),
		middleware.Handle(controller.RemoveItemTypeProperty),
	))

	return mux
}
