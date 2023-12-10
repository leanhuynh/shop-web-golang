package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(SessionLoad)
	mux.Use(InitCart)

	mux.Get("/", app.HomePage)
	// mux.Get("/checkout", app.CheckOut)
	// mux.Get("/contact", app.Contact)
	// mux.Get("/login", app.LogIn)
	// mux.Get("/my-account", app.ShowAccount)
	// mux.Route("/product", func(r chi.Router) {
	// 	mux.Get("/{productId}", app.ShowProduct)
	// })

	// mux.Route("/users", func(mux chi.Router) {
	// 	mux.Get("/checkout", app.CheckOut)
	// 	mux.Post("/placeorders", app.PlaceOrders)
	// })

	mux.Route("/products", func(mux chi.Router) {
		mux.Get("/", app.ShowProduct)
		mux.Get("/cart", app.ShowCart)
		mux.Get("/{product_id}", app.ShowProductDetail)
		mux.Post("/cart", app.AddCart)
		// mux.Post("/cart", app.AddCart)
		// mux.Put("/cart", app.UpdateCart)
		// mux.Delete("/cart", app.RemoveProduct)
		// mux.Delete("/cart/all", app.ClearCart)
	})

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
