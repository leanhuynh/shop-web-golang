package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(SessionLoad)
	// mux.Use(InitCart)

	mux.Route("/", func(mux chi.Router) {
		mux.Get("/", app.HomePage)
		mux.Get("/register", app.Register)
		mux.Get("/login", app.Login)
		mux.Get("/forgot-password", app.ForgotPassword)
		mux.Get("/reset-password", app.ResetPassword)
		mux.Get("/contact", app.Contact)
	})

	mux.Route("/products", func(mux chi.Router) {
		mux.Get("/", app.ShowProduct)
		mux.Get("/cart", app.ShowCart)
		mux.Get("/{product_id}", app.ShowProductDetail)
	})

	mux.Route("/users", func(mux chi.Router) {
		mux.Get("/checkout", app.CheckOut)
		// mux.Get("/checkout/card", app.ChargeCard)
		mux.Get("/profile", app.Profile)
	})

	mux.NotFound(app.ErrorPage)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
