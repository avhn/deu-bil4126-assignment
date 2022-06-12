package main

import (
	"ebarter/barter"
	barterdb "ebarter/barter/db"
	"ebarter/inventory"
	inventorydb "ebarter/inventory/db"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	defer barterdb.CloseDB()
	defer inventorydb.CloseDB()
	http.ListenAndServe(":8080", router())
}

func router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// barter
	r.Group(func(r chi.Router) {
		r.Post("/signup", barter.Signup)
		r.Post("/order", barter.Order)
	})

	// inventory
	r.Route("/inventory", func(r chi.Router) {
		r.Post("/add", inventory.Add)
		r.Delete("/del", inventory.Del)
		r.Put("/update", inventory.Update)
		r.Get("/list", inventory.ListAll)
		r.Get("/check", inventory.Check)
		r.Get("/cost", inventory.Cost)
		r.Get("/calculate", inventory.Calculate)
	})

	return r
}
