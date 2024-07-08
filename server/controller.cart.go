package main

import (
	"fmt"
	"myapp/internal/models"
	"net/http"
)

// hàm lấy thông tin giỏ hàng của người dùng
func (app *application) GetCartInfo(w http.ResponseWriter, r *http.Request) {

	// lấy session
	session, err := app.Redis.Get(r, app.config.session_key)
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	data := make(map[string]interface{})

	/*
	 * lấy email của người dùng
	 * nếu người dùng đã đăng nhập thì lấy email trong session
	 * nếu người dùng chưa đăng nhập hoặc giá trị lưu trong session bị lỗi --> userEmail = ""
	 */
	var userEmail string
	if email, ok := session.Values["UserEmail"].(string); ok { // nếu người dùng đa đăng nhập
		userEmail = email
	} else {
		app.errorLog.Println("Người dùng chưa đăng nhập")
		data["Message"] = "Người dùng chưa đăng nhập"
		if err := app.renderTemplate(w, r, "cart", &templateData{Data: data}); err != nil {
			app.errorLog.Println(err.Error())
			panic(err)
		}
		return
	}

	// lấy thông tin của người dùng
	user, err := app.DB.GetUserByEmail(userEmail)
	if err != nil {
		app.errorLog.Println("Không tìm thấy email người dùng trong database")
		data["Message"] = "Không tìm thấy email người dùng trong database"
		if err := app.renderTemplate(w, r, "cart", &templateData{Data: data}); err != nil {
			app.errorLog.Println(err.Error())
			panic(err)
		}
		return
	}

	// lấy thông tin giỏ hàng của người dùng
	cartDetail, err := app.DB.GetCartForUser(user.ID)
	if err != nil {
		app.errorLog.Println("Không thể lấy thông tin giỏ hàng của người dùng")
		data["Message"] = "Không thể lấy thông tin giỏ hàng của người dùng"
		if err := app.renderTemplate(w, r, "cart", &templateData{Data: data}); err != nil {
			app.errorLog.Println(err.Error())
			panic(err)
		}
		return
	}

	// tính toán thông tin giỏ hàng của người dùng ()
	var cart models.Cart
	var quantity int
	var subTotal float64
	var total float64

	// khởi tạo value cho biến
	quantity = 0
	subTotal = 0
	total = 0

	for index := 0; index < len(cartDetail); index++ {
		cart.CartDetail = append(cart.CartDetail, cartDetail[index])
		quantity += cartDetail[index].Quantity
		subTotal += float64(cartDetail[index].Total)
		total += float64(cartDetail[index].Total)
	}
	cart.Quantity = quantity
	cart.SubTotal = subTotal
	cart.Total = total

	// thêm thông tin giỏ hàng vào render
	app.infoLog.Println("Lấy thông tin giỏ hàng thành công")
	data["Cart"] = cart
	fmt.Println(cart)
	// data["Message"] = "lấy thông tin giỏ hàng thành công"
	// render page
	if err := app.renderTemplate(w, r, "cart", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}
}

// hàm thêm sản phẩm vào giỏ hàng
func (app *application) AddProduct(w http.ResponseWriter, r *http.Request) {
	// chứa thông tin từ request
	var data struct {
		ProductId int `json:"id"`
		Quantity  int `json:"quantity"`
	}

	var payLoad struct {
		Message string `json:"message"`
	}

	// lấy thông tin từ session
	session, err := app.Redis.Get(r, app.config.session_key)
	if err != nil {
		app.errorLog.Println("Không thể lấy dữ liệu từ session")
		payLoad.Message = "Không thể lấy dữ liệu từ session"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// đọc thông tin từ request
	err = app.readJSON(w, r, &data)
	if err != nil {
		app.errorLog.Println("Không thể đọc dữ liệu từ request")
		payLoad.Message = "Không thể đọc dữ liệu từ request"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// lấy thông tin userId trong session
	var userEmail string
	if email, ok := session.Values["UserEmail"].(string); !ok { // nếu không có thông tin userId trong session
		app.errorLog.Println("Session không chưa email người dùng")
		payLoad.Message = "Session không chưa email người dùng"
		app.writeJSON(w, 400, payLoad)
		return
	} else {
		userEmail = email // gán giá trị user email
	}

	// lấy thông tin product
	product, err := app.DB.GetProductById(data.ProductId)
	if err != nil {
		app.errorLog.Println(err.Error())
		payLoad.Message = "Không thể tìm thấy productId trong database"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// lấy thông tin user trong database
	user, err := app.DB.GetUserByEmail(userEmail)
	if err != nil {
		app.errorLog.Println(err.Error())
		payLoad.Message = "Không tìm thấy người dùng trong database"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// thực hiện truy vấn trong cơ sở dữ liệu
	err = app.DB.AddProducToCart(user.ID, product, data.Quantity)
	if err != nil {
		app.errorLog.Println(err.Error())
		payLoad.Message = err.Error()
		app.writeJSON(w, 400, payLoad)
		return
	}

	// gửi response tới client
	app.infoLog.Println("Thêm sản phẩm vào giỏ hàng thành công")
	payLoad.Message = "Thêm sản phẳm thành công vào giỏ hàng"
	app.writeJSON(w, 200, payLoad)
}

// hàm xóa sản phẩm ra khỏi giỏ hàng
func (app *application) RemoveProduct(w http.ResponseWriter, r *http.Request) {
	// chứa thông tin từ request
	var data struct {
		ProductId int `json:"id"`
	}

	var payLoad struct {
		Message string `json:"message"`
	}

	// lấy thông tin từ session
	session, err := app.Redis.Get(r, app.config.session_key)
	if err != nil {
		app.errorLog.Println("Không thể lấy dữ liệu từ session")
		payLoad.Message = "Không thể lấy dữ liệu từ session"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// đọc thông tin từ request
	err = app.readJSON(w, r, &data)
	if err != nil {
		app.errorLog.Println("Không thể đọc dữ liệu từ request")
		payLoad.Message = "Không thể đọc dữ liệu từ request"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// lấy thông tin userId trong session
	var userEmail string
	if email, ok := session.Values["UserEmail"].(string); !ok { // nếu không có thông tin userId trong session
		app.errorLog.Println("Session không chưa email người dùng")
		payLoad.Message = "Session không chưa email người dùng"
		app.writeJSON(w, 400, payLoad)
		return
	} else {
		userEmail = email // gán giá trị user email
	}

	// lấy thông tin user trong database
	user, err := app.DB.GetUserByEmail(userEmail)
	if err != nil {
		app.errorLog.Println(err.Error())
		payLoad.Message = "Không tìm thấy người dùng trong database"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// xóa product ra khỏi giỏ hàng
	err = app.DB.RemoveProduct(user.ID, data.ProductId)

	if err != nil {
		app.errorLog.Println(err.Error())
		payLoad.Message = err.Error()
		app.writeJSON(w, 400, payLoad)
		return
	}

	app.infoLog.Println("Xóa sản phẩm ra khỏi giỏ hàng thành công")
	payLoad.Message = "Xóa sản phẩm ra khỏi giỏ hàng thành công"
	app.writeJSON(w, 200, payLoad)
}

// hàm xóa giỏ hàng
func (app *application) RemoveCart(w http.ResponseWriter, r *http.Request) {

	var payLoad struct {
		Message string `json:"message"`
	}

	// lấy thông tin từ session
	session, err := app.Redis.Get(r, app.config.session_key)
	if err != nil {
		app.errorLog.Println("Không thể lấy dữ liệu từ session")
		payLoad.Message = "Không thể lấy dữ liệu từ session"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// lấy thông tin userId trong session
	var userEmail string
	if email, ok := session.Values["UserEmail"].(string); !ok { // nếu không có thông tin userId trong session
		app.errorLog.Println("Session không chưa email người dùng")
		payLoad.Message = "Session không chưa email người dùng"
		app.writeJSON(w, 400, payLoad)
		return
	} else {
		userEmail = email // gán giá trị user email
	}

	// lấy thông tin user trong database
	user, err := app.DB.GetUserByEmail(userEmail)
	if err != nil {
		app.errorLog.Println(err.Error())
		payLoad.Message = "Không tìm thấy người dùng trong database"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// xóa giỏ hàng
	err = app.DB.RemoveCart(user.ID)

	if err != nil {
		app.errorLog.Println(err.Error())
		payLoad.Message = err.Error()
		app.writeJSON(w, 400, payLoad)
		return
	}

	app.infoLog.Println("Xóa giỏ hàng thành công")
	payLoad.Message = "Xóa giỏ hàng thành công"
	app.writeJSON(w, 200, payLoad)
}

// hàm cập nhật số lượng sản phẩm trong giỏ hàng
func (app *application) UpdateCart(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ProductId int `json:"id"`
		Quantity  int `json:"quantity"`
	}

	var payLoad struct {
		Message string `json:"message"`
	}

	// đọc data
	err := app.readJSON(w, r, &data)
	if err != nil {
		app.errorLog.Println("Không thể đọc thông tin từ request")
		payLoad.Message = "Không thể đọc thông tin từ request"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// lấy session
	session, err := app.Redis.Get(r, app.config.session_key)
	if err != nil {
		app.errorLog.Println("Không thể lấy thông tin từ session")
		payLoad.Message = "Không thể lấy thông tin từ session"
		app.writeJSON(w, 400, payLoad)
		return
	}

	var userEmail string
	if email, ok := session.Values["UserEmail"].(string); ok { // neeus to
		userEmail = email
	} else {
		app.errorLog.Println("Người dùng chưa đăng nhập")
		payLoad.Message = "Người dùng chưa đăng nhập"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// lấy thông tin người dùng
	user, err := app.DB.GetUserByEmail(userEmail)
	if err != nil {
		app.errorLog.Println("Không tìm thấy người dùng trong database")
		payLoad.Message = "Không tìm thấy người dùng trong database"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// lấy thông tin sản phẩm
	product, err := app.DB.GetProductById(data.ProductId)
	if err != nil {
		app.errorLog.Println("Không thể lấy thông tin sản phẩm")
		payLoad.Message = "Không thể lấy thông tin sản phẩm"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// check if bonous quantity is not larger then quantity of product
	if data.Quantity > product.Quantity {
		app.errorLog.Println("Vượt quá số lượng hiện có")
		payLoad.Message = "Vượt quá số lượng hiện có"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// thực thi truy vấn
	err = app.DB.UpdateProductToCart(user.ID, product, data.Quantity)
	if err != nil {
		app.errorLog.Println(err.Error())
		payLoad.Message = err.Error()
		app.writeJSON(w, 400, payLoad)
		return
	}

	// gửi phản hồi
	app.infoLog.Println("Cập nhật thành công")
	payLoad.Message = "Cập nhật thành công"
	app.writeJSON(w, 200, payLoad)
}
