package server

import (
	"FakeAPI/internal/server/api"
	"github.com/go-chi/chi/v5"
)

func RouteFakeApi(r chi.Router) {
	r.Get("/sync", api.GetAllSync)
	r.Get("/async", api.GetAllAsync)
	r.Get("/products/compare", api.CompareProducts)
	r.Post("/dbuser", api.AddUserToDB)
	r.Put("/dbuser", api.UpsertUserToDB)
	r.Get("/dbuser", api.GetUserFromDB)
}
