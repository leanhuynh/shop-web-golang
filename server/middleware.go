package main

import (
	"errors"
	"fmt"
	"myapp/internal/urlsigner"
	"net/http"
)

// hàm xác thực token trong header
func (app *application) authenticateToken(r *http.Request) error {
	// lấy token
	theURL := r.RequestURI
	testURL := fmt.Sprintf("%s:%d%s", app.config.host, app.config.port, theURL)

	signer := urlsigner.Signer{
		Secret: []byte(app.config.secret_key),
	}

	// kiểm tra token
	valid := signer.VerifyToken(testURL)

	if !valid {
		return errors.New("token không hợp lệ")
	}

	// đảm bảo token không quá hạn
	expired := signer.Expired(testURL, 60)
	if expired {
		return errors.New("token quá hạn")
	}

	return nil
}

// hàm check trạng thái logging của người dùng
func (app *application) CheckLogged(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// lấy thông tin trong session
		session, err := app.Redis.Get(r, app.config.session_key)
		if err != nil {
			app.errorLog.Println(err.Error())
			panic(err)
		}

		// nếu chưa đăng nhập, đưa vể trang login để xử lý
		if isLogged, ok := session.Values["IsLogged"].(bool); ok && isLogged { // là true
			// nếu vẫn đang trạng thái đăng nhập --> không cần chuyển url
			// xử lí tiếp middleware khác
			next.ServeHTTP(w, r)
		} else { // là nil hoặc false
			session.Values["IsLogged"] = false // xóa trạng thái đăng nhập
			session.Values["UserEmail"] = nil  // xóa email của người dùng
			// chưa đăng nhập --> chuyển sang trang LogIn
			http.Redirect(w, r, fmt.Sprintf("%s:%d/login?reqURL=%s", app.config.host, app.config.port, r.URL.String()), http.StatusSeeOther)
			return
		}
	})
}

// func (app *application) CheckLoggedForCart (next http.Handler)

// hàm xác thực token trong url
func (app *application) checkToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// kiểm tra token có còn hợp lệ hay không
		err := app.authenticateToken(r)

		if err != nil { // nếu token không hợp lệ
			app.errorLog.Println(err.Error())
			panic(err)
		}
		next.ServeHTTP(w, r) // tiếp tục thực thi
	})
}

// hàm xử lí các lỗi error trong middleware
func (app *application) errorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			data := make(map[string]interface{})

			// render error page
			if err := recover(); err != nil {
				data["Error"] = err
				if err := app.renderTemplate(w, r, "error", &templateData{Data: data}); err != nil {
					app.errorLog.Println(err.Error())
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// hàm xử lý nếu không tìm thấy route đúng để xử lý
func (app *application) FileNotFound(w http.ResponseWriter, r *http.Request) {
	// lấy error message
	data := make(map[string]interface{})
	data["Error"] = "File not found"

	// render error page
	if err := app.renderTemplate(w, r, "error", &templateData{Data: data}); err != nil {
		app.errorLog.Println(err)
	}
}

// hàm xử lý nếu không tìm thấy route + method đúng để xử lý
func (app *application) MethodNotFound(w http.ResponseWriter, r *http.Request) {
	// lấy error message
	data := make(map[string]interface{})
	data["Error"] = "Method Not Found"

	// render error page
	if err := app.renderTemplate(w, r, "error", &templateData{Data: data}); err != nil {
		app.errorLog.Println(err)
	}
}
