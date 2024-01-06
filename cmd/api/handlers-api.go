package main

import (
	"errors"
	"fmt"
	"myapp/internal/models"
	"myapp/internal/urlsigner"
	"net/http"
	"strings"
	"time"
)

func (app *application) Register(w http.ResponseWriter, r *http.Request) {
	var data struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Mobile    string `json:"mobile"`
	}

	err := app.readJSON(w, r, &data)
	var payload struct {
		Message string `json:"message"`
	}

	if err != nil {
		app.errorLog.Println(err)
		payload.Message = "Can not read account's information"
		app.writeJSON(w, http.StatusBadRequest, payload)
		return
	}

	err = app.DB.CreateUser(data.Email, data.Password, data.FirstName, data.LastName, data.Mobile)
	if err != nil {
		payload.Message = err.Error()
	} else {
		payload.Message = "Create account successfully"
	}
	app.writeJSON(w, 200, payload)
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &data)
	var payload struct {
		Message string        `json:"message"`
		Token   *models.Token `json:"authentication_token"`
	}

	payload.Token = nil

	if err != nil {
		app.errorLog.Println(err)
		payload.Message = "Can not read account's information"
		app.writeJSON(w, 400, payload)
		return
	}

	// tim kiem thong tin nguoi dung co email tuong ung
	user, err := app.DB.GetUserByEmail(data.Email)
	if err != nil {
		payload.Message = err.Error()
		app.writeJSON(w, 400, payload)
		return
	} else {

		// compare input and database (password)
		isValid, _ := app.passwordMatches(user.Password, data.Password)

		if !isValid {
			payload.Message = "Not match password"
			app.writeJSON(w, 400, payload)
			return
		}
	}

	// generate token
	token, err := models.GenerateToken(user.ID, 2*time.Hour, models.ScopeAuthentication)
	if err != nil {
		payload.Message = "Can not generate authenticate token"
		app.writeJSON(w, 400, payload)
		return
	}

	// save token to database
	err = app.DB.InsertToken(token, user)
	if err != nil {
		payload.Message = err.Error()
		app.writeJSON(w, 400, payload)
		return
	}

	payload.Message = "Login successfully"
	payload.Token = token

	app.writeJSON(w, 200, payload)
}

func (app *application) authenticateToken(r *http.Request) (*models.User, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, errors.New("no authorization header received")
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("no authorization header received")
	}

	token := headerParts[1]
	if len(token) != 26 {
		return nil, errors.New("authentication token wrong size")
	}

	// get the user from the tokens table
	user, err := app.DB.GetUserForToken(token)
	if err != nil {
		return nil, errors.New("no matching user found")
	}

	return user, nil
}

func (app *application) CheckAuthentication(w http.ResponseWriter, r *http.Request) {
	// validate the token, and get associated user
	var payLoad struct {
		Message string `json:"message"`
	}

	// authenticate token
	user, err := app.authenticateToken(r)
	if err != nil {
		payLoad.Message = err.Error()
		app.writeJSON(w, 400, payLoad)
		return
	}

	// gui response to client neu xac thuc thanh cong
	payLoad.Message = fmt.Sprintf("authenticated user %s", user.Email)
	app.writeJSON(w, 200, payLoad)
}

func (app *application) SendPasswordResetEmail(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Email string `json:"email"`
	}

	var payload struct {
		Message string `json:"message"`
	}

	err := app.readJSON(w, r, &userInput)
	if err != nil {
		payload.Message = "Can not read account's information"
		app.writeJSON(w, 400, payload)
		return
	}

	// verify that email exists
	_, err = app.DB.GetUserByEmail(userInput.Email)

	if err != nil {
		payload.Message = "No matching user with given email"
		app.writeJSON(w, 400, payload)
		return
	}

	link := fmt.Sprintf("%s/reset-password?email=%s", app.config.frontend, userInput.Email)

	sign := urlsigner.Signer{
		Secret: []byte(app.config.secretkey),
	}

	signedLink := sign.GenerateTokenFromString(link)

	var data struct {
		Link string
	}

	data.Link = signedLink

	// send mail
	err = app.SendMail("info@widgets.com", userInput.Email, "Password Reset Request", "password-reset", data)
	if err != nil {
		payload.Message = "Error when sending reset email"
		app.writeJSON(w, 400, payload)
		return
	}

	payload.Message = "Send reset email successfully"
	app.writeJSON(w, 200, payload)
}

func (app *application) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var payload struct {
		Message string `json:"message"`
	}

	err := app.readJSON(w, r, &userInput)

	if err != nil {
		payload.Message = "Can not read input information"
		app.writeJSON(w, 400, payload)
		return
	}

	// encyrptor := encryption.Encryption{
	// 	Key: []byte(app.config.secretkey),
	// }

	// realEmail, err := encyrptor.Decrypt(userInput.Email)
	// if err != nil {
	// 	payload.Message = err.Error()
	// 	app.writeJSON(w, 400, payload)
	// 	return
	// }

	_, err = app.DB.GetUserByEmail(userInput.Email)
	if err != nil {
		payload.Message = "Can not load user with given email"
		app.writeJSON(w, 400, payload)
		return
	}

	err = app.DB.UpdatePasswordForUser(userInput.Email, userInput.Password)
	if err != nil {
		payload.Message = "Update password process failed"
		app.writeJSON(w, 400, payload)
		return
	}

	payload.Message = "Update password successfully"
	app.writeJSON(w, 200, payload)
}

// ham cho phep xoa token trong database
func (app *application) LogOut(w http.ResponseWriter, r *http.Request) {
	var payLoad struct {
		Message string `json:"message"`
	}

	var userInput struct {
		Token string `json:"token"`
	}

	err := app.readJSON(w, r, &userInput)
	if err != nil {
		payLoad.Message = "Can not read token"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// delete token from database
	// _, err = app.DB.DeleteUserForToken(userInput.Token)

	// can not delete
	// if err != nil {
	// 	payLoad.Message = err.Error()
	// 	app.writeJSON(w, 400, payLoad)
	// 	return
	// }

	// delete thanh cong
	payLoad.Message = "Log out successfully"
	app.writeJSON(w, 200, payLoad)
}

// ham lay user_id va tra ve address
func (app *application) GetUserAddress(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("Authorization")
	var payLoad struct {
		Message string           `json:"message"`
		Address []models.Address `json:"address"`
	}

	payLoad.Address = nil

	headerParts := strings.Split(authorizationHeader, " ")

	token := headerParts[1]

	// get the user from the tokens table
	user, _ := app.DB.GetUserForToken(token)
	address, err := app.DB.GetAddressForUser(user.ID)
	if err != nil {
		payLoad.Message = err.Error()
		app.writeJSON(w, 400, payLoad)
		return
	}

	payLoad.Message = "Get address successfully"
	payLoad.Address = address
	fmt.Println(address)
	app.writeJSON(w, 200, payLoad)
}

// ham xac nhan order
func (app *application) Order(w http.ResponseWriter, r *http.Request) {
	var infor struct {
		Token    string  `json:"token"`
		Quantity int     `json:"quantity"`
		Price    float32 `json:"price"`
	}

	var payLoad struct {
		Message    string `json:"message"`
		PlaceOrder string `json:"place_order"`
	}

	// read data
	err := app.readJSON(w, r, &infor)
	if err != nil {
		payLoad.Message = "Can not read data"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// insert order into database
	_, err = app.DB.InsertOrder(infor.Token, infor.Quantity, infor.Price)
	if err != nil {
		payLoad.Message = err.Error()
		app.writeJSON(w, 400, payLoad)
		return
	}

	// get user for token
	user, err := app.DB.GetUserForToken(infor.Token)
	if err != nil {
		payLoad.Message = err.Error()
		app.writeJSON(w, 400, payLoad)
		return
	}

	// send place order to user's email
	var data struct {
		Quantity  int     `json:"quantity"`
		Price     float32 `json:"price"`
		FirstName string  `json:"first_name"`
		LastName  string  `json:"last_name"`
	}
	data.Quantity = infor.Quantity
	data.Price = infor.Price
	data.FirstName = user.FirstName
	data.LastName = user.LastName

	err = app.SendMail("info@widgets.com", user.Email, "PlaceOrder Information", "send-placeorder", data)
	if err != nil {
		payLoad.Message = "Error when sending placeorder"
		app.writeJSON(w, 400, payLoad)
		return
	}

	payLoad.Message = "Order successfully"
	app.writeJSON(w, 200, payLoad)
}

// // ham lay thong tin cart trong session
// func (app *application) GetCartInfo(w http.ResponseWriter, r *http.Request) {
// 	cart := session.Get(r.Context(), "cart").(Cart)
// 	var payLoad struct {
// 		Message  string `json:"message"`
// 		CartInfo Cart   `json:"cart_info"`
// 	}

// 	// send request to client
// 	payLoad.Message = "Get cart information successfully"
// 	payLoad.CartInfo = cart.GetCart()
// 	app.writeJSON(w, 200, payLoad)
// }

// // ham them san pham vao cart
// func (app *application) AddCart(w http.ResponseWriter, r *http.Request) {
// 	var infor struct {
// 		ID       int `json:"id"`
// 		Quantity int `json:"quantity"`
// 	}

// 	var payLoad struct {
// 		Message string `json:"message"`
// 	}

// 	// read data
// 	err := app.readJSON(w, r, &infor)
// 	if err != nil {
// 		payLoad.Message = "Can not read information"
// 		app.writeJSON(w, 400, payLoad)
// 		return
// 	}

// 	cart := session.Get(r.Context(), "cart").(Cart)
// 	fmt.Println(cart)
// 	// session.Put(r.Context(), "test", &Cart{Quantity: 10})

// 	// query product information
// 	product, err := app.DB.GetProductById(infor.ID)
// 	if err != nil {
// 		payLoad.Message = "Can not find product information by id"
// 		app.writeJSON(w, 400, payLoad)
// 		return
// 	}

// 	// check if bonous quantity is not larger then quantity of product
// 	if infor.Quantity > product.Quantity {
// 		payLoad.Message = "San pham them vao vuot qua so luong hien co"
// 		app.writeJSON(w, 400, payLoad)
// 		return
// 	}
// 	// check if cart exist product_id or not,
// 	// if yes, return index of product_id

// 	index := cart.GetIndexOfProductId(infor.ID)

// 	// add product information to cart
// 	cart.AddCart(product, infor.Quantity, index)
// 	session.Put(r.Context(), "cart", cart.GetCart())
// 	cart = session.Get(r.Context(), "cart").(Cart)
// 	fmt.Println(cart)

// 	// send response to client
// 	payLoad.Message = "Add to cart successfully"
// 	app.writeJSON(w, 200, payLoad)
// }

func (app *application) GetCartInfo(w http.ResponseWriter, r *http.Request) {
	token, _ := app.getTokenOfHeader(r)

	var payLoad struct {
		Message string              `json:"message"`
		Cart    []models.CartDetail `json:"cart"`
	}

	// query for user information
	user, err := app.DB.GetUserForToken(token)
	if err != nil {
		payLoad.Message = "Can not find user"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// send cart detail
	CartDetailList, err := app.DB.GetCartForUser(user.ID)
	if err != nil {
		payLoad.Message = err.Error()
		app.writeJSON(w, 400, payLoad)
		return
	}

	// tin tong so luong product
	// quantity := 0
	// for _, item := range CartDetailList {
	// 	quantity += item.Quantity
	// }

	// send response
	payLoad.Cart = CartDetailList
	payLoad.Message = "Get cart's information successfully"
	app.writeJSON(w, 200, payLoad)
}

func (app *application) AddProduct(w http.ResponseWriter, r *http.Request) {
	var infor struct {
		ProductId int `json:"id"`
		Quantity  int `json:"quantity"`
	}

	var payLoad struct {
		Message string `json:"message"`
	}

	// read data
	err := app.readJSON(w, r, &infor)
	if err != nil {
		payLoad.Message = "Can not read information"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// get token of header
	token, _ := app.getTokenOfHeader(r)

	// query user's information
	user, _ := app.DB.GetUserForToken(token)

	// query product information
	product, err := app.DB.GetProductById(infor.ProductId)
	if err != nil {
		payLoad.Message = "Can not find product information by id"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// query
	err = app.DB.AddProducToCart(user.ID, product, infor.Quantity)
	if err != nil {
		payLoad.Message = err.Error()
		app.writeJSON(w, 400, payLoad)
		return
	}

	// send response to client
	payLoad.Message = "Add to cart successfully"
	app.writeJSON(w, 200, payLoad)
}

func (app *application) RemoveProduct(w http.ResponseWriter, r *http.Request) {
	var infor struct {
		ProductId int `json:"id"`
	}

	var payLoad struct {
		Message string `json:"message"`
	}

	// read data
	err := app.readJSON(w, r, &infor)
	if err != nil {
		payLoad.Message = "Can not read information"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// get token of header
	token, _ := app.getTokenOfHeader(r)

	// query user's information
	user, _ := app.DB.GetUserForToken(token)

	// remove product
	err = app.DB.RemoveProduct(user.ID, infor.ProductId)

	if err != nil {
		payLoad.Message = err.Error()
		app.writeJSON(w, 400, payLoad)
		return
	}

	payLoad.Message = "Remove product successfully"
	app.writeJSON(w, 200, payLoad)
}

func (app *application) RemoveCart(w http.ResponseWriter, r *http.Request) {
	var payLoad struct {
		Message string `json:"message"`
	}

	// get token of header
	token, _ := app.getTokenOfHeader(r)

	// query user's information
	user, _ := app.DB.GetUserForToken(token)

	// remove product
	err := app.DB.RemoveCart(user.ID)

	if err != nil {
		payLoad.Message = err.Error()
		app.writeJSON(w, 400, payLoad)
		return
	}

	payLoad.Message = "Remove cart successfully"
	app.writeJSON(w, 200, payLoad)
}

func (app *application) UpdateCart(w http.ResponseWriter, r *http.Request) {
	var infor struct {
		ProductId int `json:"id"`
		Quantity  int `json:"quantity"`
	}

	var payLoad struct {
		Message string `json:"message"`
	}

	// read data
	err := app.readJSON(w, r, &infor)
	if err != nil {
		payLoad.Message = "Can not read information"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// get token of header
	token, _ := app.getTokenOfHeader(r)

	// query user's information
	user, _ := app.DB.GetUserForToken(token)

	// query product information
	product, err := app.DB.GetProductById(infor.ProductId)
	if err != nil {
		payLoad.Message = "Can not find product information by id"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// check if bonous quantity is not larger then quantity of product
	if infor.Quantity > product.Quantity {
		payLoad.Message = "Vuot qua so luong hien co"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// update
	err = app.DB.UpdateProductToCart(user.ID, product, infor.Quantity)
	if err != nil {
		payLoad.Message = err.Error()
		app.writeJSON(w, 400, payLoad)
		return
	}

	// send response
	payLoad.Message = "Update thanh cong"
	app.writeJSON(w, 200, payLoad)
}
