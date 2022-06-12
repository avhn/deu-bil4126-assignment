package main

import (
	"ebarter/art"
	adb "ebarter/art/db"
	"ebarter/barter"
	barterdb "ebarter/barter/db"
	"ebarter/construction"
	cidb "ebarter/construction/db"
	"ebarter/electronics"
	edb "ebarter/electronics/db"
	"ebarter/food"
	fdb "ebarter/food/db"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func deferall() {
	barterdb.CloseDB()
	cidb.CloseDB()
	edb.CloseDB()
	adb.CloseDB()
	fdb.CloseDB()
}

func main() {
	defer deferall()
	http.ListenAndServe(":8080", router())
}

// routes
func router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// barter
	r.Group(func(r chi.Router) {
		r.Post("/signup", barter.Signup)
		r.Post("/order", barter.Order)
	})

	// art inventory
	r.Route("/art", func(r chi.Router) {
		r.Post("/add", art.Add)
		r.Delete("/del", art.Del)
		r.Put("/update", art.Update)
		r.Get("/list", art.ListAll)
		r.Get("/check", art.Check)
		r.Get("/cost", art.Cost)
		r.Get("/calculate", art.Calculate)
	})
	// construction inventory
	r.Route("/construction", func(r chi.Router) {
		r.Post("/add", construction.Add)
		r.Delete("/del", construction.Del)
		r.Put("/update", construction.Update)
		r.Get("/list", construction.ListAll)
		r.Get("/check", construction.Check)
		r.Get("/cost", construction.Cost)
		r.Get("/calculate", construction.Calculate)
	})

	// electronics inventory
	r.Route("/electronics", func(r chi.Router) {
		r.Post("/add", electronics.Add)
		r.Delete("/del", electronics.Del)
		r.Put("/update", electronics.Update)
		r.Get("/list", electronics.ListAll)
		r.Get("/check", electronics.Check)
		r.Get("/cost", electronics.Cost)
		r.Get("/calculate", electronics.Calculate)
	})

	// food inventory
	r.Route("/food", func(r chi.Router) {
		r.Post("/add", food.Add)
		r.Delete("/del", food.Del)
		r.Put("/update", food.Update)
		r.Get("/list", food.ListAll)
		r.Get("/check", food.Check)
		r.Get("/cost", food.Cost)
		r.Get("/calculate", food.Calculate)
	})

	return r
}
