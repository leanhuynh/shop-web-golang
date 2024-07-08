package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// error handle middleware
	mux.Use(app.errorMiddleware)

	// xử lý yêu cầu xác thực
	mux.Post("/register", app.Register)
	mux.Post("/login", app.Login)
	mux.Post("/forgot-password", app.SendPasswordResetEmail)
	// không thể dùng checkToken middleware để xác thực bởi @ trong html sẽ tự convert sang %40
	// cần dùng form value
	mux.Post("/update-password", app.UpdatePassword)
	mux.Route("/reset-password", func(mux chi.Router) {
		mux.Use(app.checkToken) // middleware (xác thực token)
		mux.Get("/", app.ResetPasswordPage)
	})

	// show page
	mux.Route("/", func(mux chi.Router) {
		mux.Get("/", app.HomePage)
		mux.Get("/register", app.RegisterPage)
		mux.Get("/login", app.LoginPage)
		mux.Get("/forgot-password", app.ForgotPasswordPage)
		mux.Get("/contact", app.ContactPage)
	})

	// xử lý request liên quan tới users
	mux.Route("/user", func(mux chi.Router) {
		/*
		 * kiểm tra nếu người dùng đã đăng nhập hay chưa (middleware)
		 */
		mux.Use(app.CheckLogged)

		mux.Get("/checkout", app.CheckOutPage)
		mux.Get("/profile", app.ProfilePage)
		mux.Get("/logout", app.LogOut)
		mux.Post("/order", app.Order)
	})

	// xử lí request liên quan tới cart
	mux.Route("/cart", func(mux chi.Router) {
		/*
		 * kiểm tra nếu người dùng đã đăng nhập hay chưa (middleware)
		 */
		// mux.Use(app.CheckLogged)

		mux.Get("/", app.GetCartInfo)
		mux.Post("/add", app.AddProduct)
		mux.Post("/update", app.UpdateCart)
		mux.Delete("/remove", app.RemoveProduct)
		mux.Delete("/delete", app.RemoveCart)
	})

	// xử lý request liên quan tới products
	mux.Route("/product", func(mux chi.Router) {
		mux.Get("/", app.ProductPage)
		mux.Get("/{product_id}", app.ProductDetailPage)
	})

	// không thấy route --> trả về error page
	mux.NotFound(app.FileNotFound)

	// không thấy method + route --> trả về error page
	mux.MethodNotAllowed(app.MethodNotFound)

	// import static file
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
