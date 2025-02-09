package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"}, // in production change to frontend domain
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/users", app.AllUsers) // my test route removed in production
	r.Post("/users/register", app.RegisterUser)
	r.Post("/users/login", app.LoginUser)

	r.Route("/", func(r chi.Router) {
		r.Use(app.Authenticate)
		r.Post("/users/logout", app.LogoutUser)
		r.Get("/priorities", app.AllPriorities)

		r.Route("/todo", func(r chi.Router) {
			r.Get("/", app.AllTodos)
			r.Post("/save", app.SaveTodo)
			r.Post("/delete", app.DeleteTodo)
		})
	})

	return r
}
