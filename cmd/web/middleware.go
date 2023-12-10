package main

import (
	"encoding/gob"
	"net/http"
)

func SessionLoad(next http.Handler) http.Handler {
	// return session.LoadAndSave(next)
	return session.LoadAndSave(next)
}

func InitCart(next http.Handler) http.Handler {
	gob.Register(Cart{})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// if session.Exists(r.Context(), "cart", Cart{
		// 	CouponCode:      "code 1",
		// 	ShippingAddress: "address 1",
		// }).

		if !session.Exists(r.Context(), "cart") {
			session.Put(r.Context(), "cart", Cart{
				CouponCode:      "code 1",
				ShippingAddress: "address 1",
			})
		}

		next.ServeHTTP(w, r)
	})
}
