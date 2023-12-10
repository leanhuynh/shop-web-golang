package main

// type JsonResponse struct {
// 	OK      bool   `json:"ok"`
// 	Message string `json:"message,omitempty"`
// 	Content string `json:"content,omitempty"`
// 	ID      int    `json:"id,omitempty"`
// }

// func (app *application) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
// 	email := chi.URLParam(r, "email")

// 	u, err := app.DB.GetUserByEmail(email)

// 	if err != nil {
// 		app.errorLog.Println(err)
// 		return
// 	}

// 	fmt.Println(u)

// 	out, err := json.MarshalIndent(u, "", "   ")
// 	if err != nil {
// 		app.errorLog.Println(err)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(out)
// }

// func (app *application) CreateAuthToken(w http.ResponseWriter, r *http.Request) {
// 	var userInput struct {
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}

// 	err := r.ParseForm()
// 	if err != nil {
// 		app.badRequest(w, r, err)
// 		return
// 	}

// 	userInput.Email = string(r.Form.Get("username"))
// 	userInput.Password = string(r.Form.Get("password"))

// 	// get the user from the database by email; send error if invalid email
// 	user, err := app.DB.GetUserByEmail(userInput.Email)
// 	if err != nil {
// 		app.invalidCredentials(w)
// 		return
// 	}

// 	// validate the password; send error if invalid password
// 	validPassword, err := app.passwordMatches(user.Password, userInput.Password)
// 	if err != nil {
// 		app.invalidCredentials(w)
// 		return
// 	}

// 	if !validPassword {
// 		app.invalidCredentials(w)
// 		return
// 	}

// 	// generate the token
// 	token, err := models.GenerateToken(user.ID, 24*time.Hour, models.ScopeAuthentication)
// 	if err != nil {
// 		app.badRequest(w, r, err)
// 		return
// 	}

// 	// save to database
// 	err = app.DB.InsertToken(token, user)
// 	if err != nil {
// 		app.badRequest(w, r, err)
// 		return
// 	}

// 	// send response

// 	var payload struct {
// 		Error   bool          `json:"error"`
// 		Message string        `json:"message"`
// 		Token   *models.Token `json:"authentication_token"`
// 	}
// 	payload.Error = false
// 	payload.Message = fmt.Sprintf("token for %s created", userInput.Email)
// 	payload.Token = token

// 	_ = app.writeJSON(w, http.StatusOK, payload)
// }

// func (app *application) authenticationToken(r *http.Request) (*models.User, error) {
// 	authorizationHeader := r.Header.Get("Authorization")
// 	fmt.Println(authorizationHeader)
// 	if authorizationHeader == "" {
// 		return nil, errors.New("no authorization header received")
// 	}

// 	headerParts := strings.Split(authorizationHeader, " ")
// 	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
// 		return nil, errors.New("no authorization header received")
// 	}

// 	token := headerParts[1]
// 	if len(token) != 26 {
// 		return nil, errors.New("authentication token wrong size")
// 	}

// 	// get the user from the tokens table
// 	user, err := app.DB.GetUserForToken(token)
// 	if err != nil {
// 		return nil, errors.New("no matching user found")
// 	}

// 	return user, nil
// }

// func (app *application) CheckAuthentication(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("called")
// 	// validate the token, and get associated user
// 	user, err := app.authenticationToken(r)
// 	if err != nil {
// 		app.invalidCredentials(w)
// 		return
// 	}

// 	// valid user
// 	var payload struct {
// 		Error   bool   `json:"error"`
// 		Message string `json:"message"`
// 	}
// 	payload.Error = false
// 	payload.Message = fmt.Sprintf("authenticated user %s", user.Email)
// 	app.writeJSON(w, http.StatusOK, payload)
// }
