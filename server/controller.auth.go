package main

import (
	"encoding/json"
	"fmt"
	"myapp/internal/urlsigner"
	"net/http"
)

// hàm xử lí Register, trả vể reqURL nếu thành công
func (app *application) Register(w http.ResponseWriter, r *http.Request) {
	var data struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Mobile    string `json:"mobile"`
	}

	// lấy thông tin trong session
	session, err := app.Redis.Get(r, app.config.session_key)
	if err != nil {
		app.errorLog.Println(err)
		panic(err) // trả về error
	}

	Flash := make(map[string]string)

	// kiểm tra thông tin có thể parse được hay không
	if err := r.ParseForm(); err != nil {
		app.errorLog.Println(err.Error())
		panic(err) // trả về error
	}

	// thêm value từ form vào biến data
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	data.FirstName = r.FormValue("firstname")
	data.LastName = r.FormValue("lastname")
	data.Mobile = r.FormValue("mobile")

	// lấy url hiện tại mà người dùng đang đứng
	reqURL := r.FormValue("reqURL")
	if reqURL == "" {
		reqURL = "/"
	}

	// gọi hàm trong sql để xử lí
	err = app.DB.CreateUser(data.Email, data.Password, data.FirstName, data.LastName, data.Mobile)
	if err != nil { // nếu sql xảy ra lỗi
		app.errorLog.Println(err.Error())
		Flash["RegisterMessage"] = "Không thể tạo tài khoản mới"
		jsonData, _ := json.Marshal(Flash)
		session.AddFlash(string(jsonData))
		session.Save(r, w)

		// render về page login
		http.Redirect(w, r, fmt.Sprintf("%s:%d/register?reqURL=%s", app.config.host, app.config.port, reqURL), http.StatusTemporaryRedirect)
		return
	} else {
		app.infoLog.Println("Tạo tài khoản mới thành công")
		Flash["RegisterMessage"] = "Tạo tài khoản mới thành công"
		jsonData, _ := json.Marshal(Flash)
		session.AddFlash(string(jsonData))
		session.Save(r, w)

		// render về page login
		http.Redirect(w, r, fmt.Sprintf("%s:%d/login?reqURL=%s", app.config.host, app.config.port, reqURL), http.StatusSeeOther)
		return
	}
}

// hàm xử lí Login, trả về reqURL nếu thành công
func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	// chứa thông tin login
	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// lấy thông tin trong session
	session, err := app.Redis.Get(r, app.config.session_key)
	if err != nil {
		app.errorLog.Println(err)
		panic(err) // trả về error
	}

	Flash := make(map[string]string)

	// kiểm tra thông tin có thể parse được hay không
	if err := r.ParseForm(); err != nil {
		app.errorLog.Println(err.Error())
		panic(err) // trả về error
	}

	// thêm value từ form vào biến data
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")

	// lấy url hiện tại mà người dùng đang đứng
	reqURL := r.FormValue("reqURL")
	if reqURL == "" {
		reqURL = "/"
	}

	// tìm kiểm thông tin người dùng có email tương ứng
	user, err := app.DB.GetUserByEmail(data.Email)
	if err != nil { // bị lỗi
		app.errorLog.Println("Email không tồn tại")
		Flash["LoginMessage"] = "Email không tồn tại"
		jsonData, _ := json.Marshal(Flash)
		session.AddFlash(string(jsonData))
		session.Save(r, w)

		// tiếp tục login page
		http.Redirect(w, r, fmt.Sprintf("%s:%d/login?reqURL=%s", app.config.host, app.config.port, reqURL), http.StatusFound)
		return
	} else {
		// so sánh input password và password trong database
		isValid, _ := app.passwordMatches(user.Password, data.Password)

		if !isValid {
			app.errorLog.Println("Password Không đúng")
			Flash["LoginMessage"] = "Password không đúng"
			jsonData, _ := json.Marshal(Flash)
			session.AddFlash(string(jsonData))
			session.Save(r, w)

			// tiếp tục login page
			http.Redirect(w, r, fmt.Sprintf("%s:%d/login?reqURL=%s", app.config.host, app.config.port, reqURL), http.StatusFound)
			return
		}
	}

	// đăng nhập thành công
	session.Values["IsLogged"] = true        // lưu trạng thái đăng nhập của người dùng
	session.Values["UserEmail"] = data.Email // lưu email của người dùng (nếu đã đăng nhập)
	app.infoLog.Println("Đăng nhập thành công")
	Flash["LoginMessage"] = "Đăng nhập thành công"
	jsonData, _ := json.Marshal(Flash)
	session.AddFlash(string(jsonData))
	session.Save(r, w)
	http.Redirect(w, r, reqURL, http.StatusSeeOther)
}

// hàm xử lí send reset-email
func (app *application) SendPasswordResetEmail(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string `json:"email"`
	}

	session, err := app.Redis.Get(r, app.config.session_key)
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	// kiểm tra có parse đươc thông tin trong body request hay không
	if err = r.ParseForm(); err != nil {
		app.errorLog.Println(err.Error())
		panic(err) // trả về error
	}

	// lấy thông tin reset-email
	data.Email = r.FormValue("email")

	Flash := make(map[string]string)

	// kiểm tra mail có tổn tại trong database
	_, err = app.DB.GetUserByEmail(data.Email)
	if err != nil {
		app.errorLog.Println(err.Error())
		Flash["SendResetMailMessage"] = "Không tìm thấy email trong database"
		jsonData, _ := json.Marshal(Flash)
		session.AddFlash(string(jsonData))
		session.Save(r, w)

		http.Redirect(w, r, fmt.Sprintf("%s:%d/forgot-password", app.config.host, app.config.port), http.StatusFound)
		return
	}

	link := fmt.Sprintf("%s:%d/reset-password?email=%s", app.config.host, app.config.port, data.Email)

	sign := urlsigner.Signer{
		Secret: []byte(app.config.secret_key),
	}

	// tạo token từ link
	signedLink := sign.GenerateTokenFromString(link)

	var load struct {
		Link string
	}

	load.Link = signedLink

	// send mail
	err = app.SendMail("info@widgets.com", data.Email, "Password Reset Request", "password-reset", load)
	if err != nil {
		app.errorLog.Println(err.Error())
		Flash["SendResetMailMessage"] = "Xảy ra lỗi trong khi gửi reset-email"
		jsonData, _ := json.Marshal(Flash)
		session.AddFlash(string(jsonData))
		session.Save(r, w)

		http.Redirect(w, r, fmt.Sprintf("%s:%d/forgot-password", app.config.host, app.config.port), http.StatusFound)
		return
	} else {
		app.infoLog.Println("Gửi reset-email thành công")
		Flash["SendResetMailMessage"] = "Gửi reset-email thành công"
		jsonData, _ := json.Marshal(Flash)
		session.AddFlash(string(jsonData))
		session.Save(r, w)

		// gửi email thành công --> trở về trang chủ
		// người dùng sẽ vào email nhận được và nhấp vào link reset-password
		http.Redirect(w, r, fmt.Sprintf("%s:%d/forgot-password", app.config.host, app.config.port), http.StatusFound)
		return
	}
}

/*
 * xử lí reset password
 * nếu thành công --> trả về trang chủ
 */
func (app *application) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// lấy thông tin trong session
	session, err := app.Redis.Get(r, app.config.session_key)
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	// kiểm tra thông tin có thể parse được hay không
	if err := r.ParseForm(); err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	// kiểm tra tính xác thực của token
	URI := r.FormValue("uri")
	testURL := fmt.Sprintf("%s:%d%s", app.config.host, app.config.port, URI)
	signer := urlsigner.Signer{
		Secret: []byte(app.config.secret_key),
	}

	// kiểm tra token
	valid := signer.VerifyToken(testURL)

	// kiểm tra tính hợp lệ
	if !valid {
		app.errorLog.Println("token không hợp lệ")
		panic("Token không hợp lệ")
	}

	// đảm bảo token không quá hạn
	expired := signer.Expired(testURL, 60)
	if expired {
		app.errorLog.Println("Token quá hạn")
		panic("Token quá hạn")
	}

	// lấy thông tin trong request body
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")

	Flash := make(map[string]string)

	// kiểm tra email có tồn tại trong database
	_, err = app.DB.GetUserByEmail(data.Email)
	if err != nil {
		app.errorLog.Println(err.Error())
		Flash["ResetPassword"] = "Không tìm thấy email trong database"
		jsonData, _ := json.Marshal(Flash)
		session.AddFlash(string(jsonData))
		session.Save(r, w)

		http.Redirect(w, r, r.URL.String(), http.StatusFound)
		return
	}

	// tiến hành cập nhật password trong database
	fmt.Println(data)
	err = app.DB.UpdatePasswordForUser(data.Email, data.Password)
	if err != nil {
		app.errorLog.Println(err.Error())
		panic("Xử lí sql bị lỗi")
	}

	// nếu cập nhật password thành công
	app.infoLog.Println("Cập nhật password thành công")
	Flash["ResetPassword"] = "Cập nhật password thành công"
	jsonData, _ := json.Marshal(Flash)
	session.AddFlash(string(jsonData))
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
}

// ham cho phep xoa token trong database
func (app *application) LogOut(w http.ResponseWriter, r *http.Request) {
	// lấy thông tin session
	session, err := app.Redis.Get(r, app.config.session_key)
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	// thay đổi trạng thái Logging của người dùng
	session.Values["IsLogged"] = false

	// xóa session
	session.Options.MaxAge = -1

	// lưu lại session
	session.Save(r, w)

	// trở về trang chủ
	http.Redirect(w, r, fmt.Sprintf("%s:%d", app.config.host, app.config.port), http.StatusSeeOther)
}
