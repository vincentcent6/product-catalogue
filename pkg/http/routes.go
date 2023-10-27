package http

import (
	"github.com/go-chi/chi/v5"
)

func NewRoutes() *chi.Mux {
	router := chi.NewRouter()

	// Product
	router.Post("/product", CreateProduct)

	return router
}
