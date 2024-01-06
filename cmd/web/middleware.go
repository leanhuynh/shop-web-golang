package main

import (
	"net/http"
)

func (app *application) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check authenticate with token
		_, err := app.authenticateToken(r)

		if err != nil {
			// send bad response to client
			var payLoad struct {
				Message string `json:"message"`
			}
			payLoad.Message = "Invalid token"
			app.writeJSON(w, 400, payLoad)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func SessionLoad(next http.Handler) http.Handler {
	// return session.LoadAndSave(next)
	return session.LoadAndSave(next)
}

// func InitCart(next http.Handler) http.Handler {
// 	gob.Register(Cart{})
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		// if session.Exists(r.Context(), "cart", Cart{
// 		// 	CouponCode:      "code 1",
// 		// 	ShippingAddress: "address 1",
// 		// }).

// 		if !session.Exists(r.Context(), "cart") {
// 			session.Put(r.Context(), "cart", Cart{
// 				CouponCode:      "code 1",
// 				ShippingAddress: "address 1",
// 				Quantity:        0,
// 			})
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }
