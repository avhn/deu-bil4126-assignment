package main

import (
	"ebarter/barter"
	barterdb "ebarter/barter/db"
	"ebarter/inventory"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	barterdb.InitSQLiteDB()
	defer barterdb.CloseDB()
	u := barterdb.NewUser("ulug@gmail.com")
	u.Create()
	fmt.Println(u)
	r := barterdb.NewRequest(u.Email, "food", "hamburger", 10, "electronics", "macbook pro", 1)
	r.Create()
	fmt.Println(r)

	http.ListenAndServe(":8080", router())

}

func router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// public
	r.Group(func(r chi.Router) {
		r.Post("/signup", barter.Signup)
	})

	r.Group(func(r chi.Router) {
		r.Post("/add", inventory.Add)
		r.Post("/add", inventory.ListAll)

	})

	return r
}
