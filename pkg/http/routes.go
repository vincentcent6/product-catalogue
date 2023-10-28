package http

import (
	"github.com/go-chi/chi/v5"
)

func NewRoutes() *chi.Mux {
	router := chi.NewRouter()

	// Product
	router.Post("/product", CreateProduct)
	router.Put("/product", UpdateProduct)
	router.Get("/product", GetProduct)
	router.Get("/products", GetProducts)

	return router
}
