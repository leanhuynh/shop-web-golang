package main

import (
	"net/http"
)

// VirtualTerminal displays the virtual terminal page
func (app *application) VirtualTerminal(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "terminal", &templateData{}, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) HomePage(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "homepage", nil); err != nil {
		app.errorLog.Println(err)
	}
}

// PaymentSucceeded displays the receipt page
func (app *application) PaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// read posted data
	// cardHolder := r.Form.Get("cardholder_name")
	// email := r.Form.Get("email")
	// paymentIntent := r.Form.Get("payment_intent")
	// paymentMethod := r.Form.Get("payment_method")
	// paymentAmount := r.Form.Get("payment_amount")
	// paymentCurrency := r.Form.Get("payment_currency")

	// card := cards.Card{
	// 	Secret: app.config.stripe.secret,
	// 	Key:    app.config.stripe.key,
	// }

	// pi, err := card.RetrievePaymentIntent(paymentIntent)
	// if err != nil {
	// 	app.errorLog.Println(err)
	// 	return
	// }

	// pm, err := card.GetPaymentMethod(paymentMethod)
	// if err != nil {
	// 	app.errorLog.Println(err)
	// 	return
	// }

	// lastFour := pm.Card.Last4
	// expiryMonth := pm.Card.ExpMonth
	// expiryYear := pm.Card.ExpYear

	// data := make(map[string]interface{})
	// data["cardholder"] = cardHolder
	// data["email"] = email
	// data["pi"] = paymentIntent
	// data["pm"] = paymentMethod
	// data["pa"] = paymentAmount
	// data["pc"] = paymentCurrency
	// data["last_four"] = lastFour
	// data["expiry_month"] = expiryMonth
	// data["expiry_year"] = expiryYear
	// data["bank_return_code"] = pi.Charges.Data[0].ID

	// should write this data to session, and then redirect user to new page?

	if err := app.renderTemplate(w, r, "succeeded", &templateData{
		// Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}
}