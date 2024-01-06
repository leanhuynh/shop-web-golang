package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	// mux.Use(SessionLoad)
	// mux.Use(InitCart)

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Put("/api/register", app.Register)

	mux.Post("/api/login", app.Login)

	mux.Post("/api/forgot-password", app.SendPasswordResetEmail)

	mux.Post("/api/reset-password", app.ResetPassword)

	// check token before render page
	mux.Post("/api/is-authenticated", app.CheckAuthentication)

	mux.Post("/api/users/logout", app.LogOut)

	// use a authentication function to check token
	mux.Route("/api/users", func(mux chi.Router) {
		mux.Use(app.Auth)

		// lay thong tin user's address
		mux.Post("/address", app.GetUserAddress)

		// them order
		mux.Put("/placeorder", app.Order)
	})

	mux.Route("/api/cart", func(mux chi.Router) {
		mux.Use(app.Auth)

		mux.Get("/", app.GetCartInfo)

		mux.Put("/add", app.AddProduct)

		mux.Post("/update", app.UpdateCart)

		mux.Delete("/remove", app.RemoveProduct)

		mux.Delete("/delete", app.RemoveCart)
	})

	return mux
}
