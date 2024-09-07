package main

import (
	"encoding/json"
	"math"
	"myapp/internal/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// hàm render trang chủ
func (app *application) HomePage(w http.ResponseWriter, r *http.Request) {
	// query for feature products
	var featureProducts []models.Product
	limitProduct := 10 // giới hạn số product được show
	featureProducts, err := app.DB.GetAllFeatureProducts(limitProduct)
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	// thêm host, port vào feature product
	var featureProductsWithHostPort []models.ProductWithHostPort
	for index := 0; index < len(featureProducts); index++ {
		var product models.ProductWithHostPort
		product.ID = featureProducts[index].ID
		product.Name = featureProducts[index].Name
		product.ImagePath = featureProducts[index].ImagePath
		product.OldPrice = featureProducts[index].OldPrice
		product.Price = featureProducts[index].Price
		product.Summary = featureProducts[index].Summary
		product.Description = featureProducts[index].Description
		product.Specification = featureProducts[index].Specification
		product.Quantity = featureProducts[index].Quantity
		product.CategoryId = featureProducts[index].CategoryId
		product.BrandId = featureProducts[index].BrandId
		product.CreatedAt = featureProducts[index].CreatedAt
		product.Stars = featureProducts[index].Stars
		product.Host = app.config.host
		product.Port = app.config.port

		featureProductsWithHostPort = append(featureProductsWithHostPort, product)
	}

	// query for category product
	categories, err := app.DB.GetAllCategories()
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	// query for recent products
	recentProducts, err := app.DB.GetRecentProducts()
	if err != nil {
		app.errorLog.Println(recentProducts)
		panic(err)
	}

	// thêm host, port vào recent products
	var recentProducsWithHostPort []models.ProductWithHostPort
	for index := 0; index < len(featureProducts); index++ {
		var product models.ProductWithHostPort
		product.ID = recentProducts[index].ID
		product.Name = recentProducts[index].Name
		product.ImagePath = recentProducts[index].ImagePath
		product.OldPrice = recentProducts[index].OldPrice
		product.Price = recentProducts[index].Price
		product.Summary = recentProducts[index].Summary
		product.Description = recentProducts[index].Description
		product.Specification = recentProducts[index].Specification
		product.Quantity = recentProducts[index].Quantity
		product.CategoryId = recentProducts[index].CategoryId
		product.BrandId = recentProducts[index].BrandId
		product.CreatedAt = recentProducts[index].CreatedAt
		product.Stars = recentProducts[index].Stars
		product.Host = app.config.host
		product.Port = app.config.port

		recentProducsWithHostPort = append(recentProducsWithHostPort, product)
	}

	//query for brands
	Brands, err := app.DB.GetAllBrands()
	if err != nil {
		app.errorLog.Println(Brands)
		panic(err)
	}

	data := make(map[string]interface{})
	data["featureProducts"] = featureProductsWithHostPort
	data["CategoryArray"] = categories
	data["RecentProducts"] = recentProducsWithHostPort
	data["Brands"] = Brands

	if err := app.renderTemplate(w, r, "index", &templateData{
		Data: data,
		Page: "HOME",
	}, "product"); err != nil {
		app.errorLog.Println(err)
		panic(err)
	}
}

// hàm render trang đang kí tài khoản
func (app *application) RegisterPage(w http.ResponseWriter, r *http.Request) {
	session, err := app.Redis.Get(r, app.config.session_key)
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	// lấy thông tin Flash message
	flashes := session.Flashes()

	// xóa flash mesage trong session
	session.Save(r, w)

	// chuyển flash message into interface[string]interface{} data types
	flashMessages := make(map[string]string)
	if len(flashes) > 0 {
		for _, flash := range flashes {
			flashMap := make(map[string]string)
			err := json.Unmarshal([]byte(flash.(string)), &flashMap)
			if err != nil {
				app.errorLog.Println(err.Error())
				panic(err)
			}
			for k, v := range flashMap {
				flashMessages[k] = v
			}
		}
	}

	// lấy reqURL trong url
	reqURL := r.URL.Query().Get("reqURL")
	if reqURL == "" {
		reqURL = "/"
	}
	data := make(map[string]interface{})
	data["ReqURL"] = reqURL

	if err := app.renderTemplate(w, r, "register", &templateData{Data: data, Flash: flashMessages}); err != nil {
		app.errorLog.Println(err)
		panic(err)
	}
}

// hàm render trang log-in
func (app *application) LoginPage(w http.ResponseWriter, r *http.Request) {
	// lấy session
	session, err := app.Redis.Get(r, app.config.session_key)
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	// lấy thông tin flash message
	flashes := session.Flashes()

	// xóa flash message trong session
	session.Save(r, w)

	// chuyển flash message into interface[string]interface{} data types
	flashMessages := make(map[string]string)
	if len(flashes) > 0 {
		for _, flash := range flashes {
			flashMap := make(map[string]string)
			err := json.Unmarshal([]byte(flash.(string)), &flashMap)
			if err != nil {
				app.errorLog.Println(err.Error())
				panic(err)
			}
			for k, v := range flashMap {
				flashMessages[k] = v
			}
		}
	}

	// lấy reqURL trong url
	reqURL := r.URL.Query().Get("reqURL")
	if reqURL == "" {
		reqURL = "/"
	}
	data := make(map[string]interface{})
	data["ReqURL"] = reqURL

	// render page Login
	if err := app.renderTemplate(w, r, "login", &templateData{Data: data, Flash: flashMessages}); err != nil {
		app.errorLog.Println(err)
		panic(err)
	}
}

// hàm render trang quên mật khẩu
func (app *application) ForgotPasswordPage(w http.ResponseWriter, r *http.Request) {
	// lấy session
	session, err := app.Redis.Get(r, app.config.session_key)
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	// lấy thông tin flash message
	flashes := session.Flashes()

	// xóa flash message trong session
	session.Save(r, w)

	// chuyển flash message into interface[string]string{} data types
	flashMessages := make(map[string]string)
	if len(flashes) > 0 {
		for _, flash := range flashes {
			flashMap := make(map[string]string)
			err := json.Unmarshal([]byte(flash.(string)), &flashMap)
			if err != nil {
				app.errorLog.Println(err.Error())
				panic(err)
			}
			for k, v := range flashMap {
				flashMessages[k] = v
			}
		}
	}

	// render forgot-password page
	if err := app.renderTemplate(w, r, "forgot-password", &templateData{Flash: flashMessages}); err != nil {
		app.errorLog.Println(err)
		panic(err)
	}
}

// hàm render trang thay đổi mật khẩu
func (app *application) ResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	// lấy session
	session, err := app.Redis.Get(r, app.config.session_key)
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	// lấy thông tin email
	email := r.URL.Query().Get("email")
	data := make(map[string]interface{})
	data["URI"] = r.RequestURI
	data["Email"] = email

	// lấy thông tin flash message
	flashes := session.Flashes()

	// xóa flash message trong session
	session.Save(r, w)

	// chuyển flash message into interface[string]interface{} data types
	flashMessages := make(map[string]string)
	if len(flashes) > 0 {
		for _, flash := range flashes {
			flashMap := make(map[string]string)
			err := json.Unmarshal([]byte(flash.(string)), &flashMap)
			if err != nil {
				app.errorLog.Println(err.Error())
				panic(err)
			}
			for k, v := range flashMap {
				flashMessages[k] = v
			}
		}
	}

	// render reset-password page
	if err := app.renderTemplate(w, r, "reset-password", &templateData{Data: data, Flash: flashMessages}); err != nil {
		app.errorLog.Println(err)
		panic(err)
	}
}

// hàm render Product List
func (app *application) ProductPage(w http.ResponseWriter, r *http.Request) {
	// lấy thông tin query trong url
	queryParams := r.URL.Query()

	// lấy thông tin sort filter
	var sort string
	if len(queryParams["sort"]) == 1 {
		sort = queryParams["sort"][0]
	}

	// nếu không nằm trong tập giá --> gán với giá trị mặc định là price
	if check := ContainString([]string{"price", "newest", "popular"}, sort); !check {
		sort = "price"
	}

	// tìm trang hiện tại đang đứng
	limit := 6
	var page int
	if len(queryParams["page"]) == 1 {
		page, _ = strconv.Atoi(queryParams["page"][0])
	}
	// nếu giá trị page không hợp lệ --> gán với giá trị mặc định là 1
	if page < 1 {
		page = 1
	}

	// lấy thông tin search-keyword
	var search string
	if len(queryParams["search"]) == 1 {
		search = queryParams["search"][0]
	}

	// lấy thông tin tagId
	var tagId int
	if len(queryParams["tag"]) == 1 {
		tagId, _ = strconv.Atoi(queryParams["tag"][0])
	}
	if tagId < 0 {
		tagId = 0
	}

	// lấy thông tin brandId
	var brandId int
	if len(queryParams["brand"]) == 1 {
		brandId, _ = strconv.Atoi(queryParams["brand"][0])
	}
	if brandId < 0 {
		brandId = 0
	}

	// lấy thông tin category
	var categoryId int
	if len(queryParams["category"]) == 1 {
		categoryId, _ = strconv.Atoi(queryParams["category"][0])
	}
	if categoryId < 0 {
		categoryId = 0
	}

	// xử lý filter products
	productList, err := app.DB.GetProductsWithFilter(sort, search, tagId, categoryId, brandId, 0, 0)
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	// tính số lượng page có thể có
	PageSize := int(math.Ceil(float64(len(productList)) / float64(limit))) // (tổng số lượng + số lượng trong 1 trang - 1) / (số lượng trong 1 trang)

	// xử lý filter products với page
	filterProducts, err := app.DB.GetProductsWithFilter(sort, search, tagId, categoryId, brandId, (page-1)*limit, limit)
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	// thêm host, port : dùng trong lúc render html
	var filterProductWithHostPorts []models.ProductWithHostPort
	for index := 0; index < len(filterProducts); index++ {
		var product models.ProductWithHostPort
		product.ID = filterProducts[index].ID
		product.Name = filterProducts[index].Name
		product.ImagePath = filterProducts[index].ImagePath
		product.OldPrice = filterProducts[index].OldPrice
		product.Price = filterProducts[index].Price
		product.Summary = filterProducts[index].Summary
		product.Description = filterProducts[index].Description
		product.Specification = filterProducts[index].Specification
		product.Quantity = filterProducts[index].Quantity
		product.CategoryId = filterProducts[index].CategoryId
		product.BrandId = filterProducts[index].BrandId
		product.CreatedAt = filterProducts[index].CreatedAt
		product.Stars = filterProducts[index].Stars
		product.Host = app.config.host
		product.Port = app.config.port

		filterProductWithHostPorts = append(filterProductWithHostPorts, product)
	}

	// xử lí originalURL
	OriginalUrl := r.URL.String()
	OriginalUrl = RemoveParams("sort", OriginalUrl)
	OriginalUrl += "?"

	// lấy thông tin về tất cả categories
	categories, err := app.DB.GetInforCategoriesWithNumofProducts()
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)

	}

	// lấy thông tin về tất cả brands
	brands, err := app.DB.GetInfoBrandsWithNumofProduct()
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	// lấy thông tin tất cả tags
	tags, err := app.DB.GetAllTags()
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	data := make(map[string]interface{})
	data["ProductList"] = filterProductWithHostPorts
	data["Sort"] = sort
	data["OriginalUrl"] = OriginalUrl
	data["Categories"] = categories
	data["Brands"] = brands
	data["Tags"] = tags
	data["CurrentPage"] = page
	data["PageSize"] = PageSize

	if err := app.renderTemplate(w, r, "product-list", &templateData{
		Data: data,
		Page: "PRODUCT",
	}, "product", "right-column"); err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}
}

// hàm render check-out page
func (app *application) CheckOutPage(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	// lấy session
	session, err := app.Redis.Get(r, app.config.session_key)
	if err != nil {
		app.errorLog.Println()
		panic(err)
	}

	// lấy email của người dùng
	var userEmail string

	if email, ok := session.Values["UserEmail"].(string); ok {
		userEmail = email
	} else {
		app.errorLog.Println("Người dùng chưa đăng nhập")
		panic("Người dùng chưa đăng nhập")
	}

	// tìm kiếm thông tin của người dùng
	user, err := app.DB.GetUserByEmail(userEmail)
	if err != nil {
		app.errorLog.Println("Không tìm thấy tài khoản trong database")
		panic("Không tìm thấy tài khoản trong database")
	}

	// lấy thông tin shipping address của người dùng
	addresses, err := app.DB.GetAddressForUser(user.ID)
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	// lấy thông tin giỏ hàng của người dùng
	cartDetail, err := app.DB.GetCartForUser(user.ID)
	if err != nil {
		app.errorLog.Println("Không thể lấy thông tin giỏ hàng của người dùng")
		panic("Không thể lấy thông tin giỏ hàng của người dùng")
	}

	// tính toán thông tin giỏ hàng của người dùng
	var cart models.Cart
	var quantity int
	var subTotal float64
	var total float64

	// khởi tạo value cho biến
	quantity = 0
	subTotal = 0
	total = 0

	// tính toán thông tin trong giỏ hàng
	for index := 0; index < len(cartDetail); index++ {
		cart.CartDetail = append(cart.CartDetail, cartDetail[index])
		quantity += cartDetail[index].Quantity
		subTotal += float64(cartDetail[index].Total)
		total += float64(cartDetail[index].Total)
	}
	cart.Quantity = quantity
	cart.SubTotal = subTotal
	cart.Total = total

	data["Addresses"] = addresses
	data["Cart"] = cart

	// render check-out page
	if err := app.renderTemplate(w, r, "checkout", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
		panic(err)
	}
}

// hàm render chi tiết sản phẩm
func (app *application) ProductDetailPage(w http.ResponseWriter, r *http.Request) {

	// lấy thông tin productId
	var productId int
	if chi.URLParam(r, "product_id") != "" {
		productId, _ = strconv.Atoi(chi.URLParam(r, "product_id"))
	}

	// lấy thông tin về tất cả categories với num-of-products
	categories, err := app.DB.GetInforCategoriesWithNumofProducts()
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	// lấy thông tin về tất cả brands với num-of-products
	brands, err := app.DB.GetInfoBrandsWithNumofProduct()
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	// lấy thông tin về tất cả tags
	tags, err := app.DB.GetAllTags()
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	// lấy thông tin về related products
	relatedProducts, err := app.DB.GetRelatedProductsInfor(productId)
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err.Error())
	}

	// thêm thông tin host, port vào relatedProducts để rendder html
	var relatedProductsWithHostPort []models.ProductWithHostPort
	for index := 0; index < len(relatedProducts); index++ {
		var product models.ProductWithHostPort
		product.ID = relatedProducts[index].ID
		product.Name = relatedProducts[index].Name
		product.ImagePath = relatedProducts[index].ImagePath
		product.OldPrice = relatedProducts[index].OldPrice
		product.Price = relatedProducts[index].Price
		product.Summary = relatedProducts[index].Summary
		product.Description = relatedProducts[index].Description
		product.Specification = relatedProducts[index].Specification
		product.Quantity = relatedProducts[index].Quantity
		product.CategoryId = relatedProducts[index].CategoryId
		product.BrandId = relatedProducts[index].BrandId
		product.CreatedAt = relatedProducts[index].CreatedAt
		product.Stars = relatedProducts[index].Stars
		product.Host = app.config.host
		product.Port = app.config.port

		relatedProductsWithHostPort = append(relatedProductsWithHostPort, product)
	}

	// lấy thông tin về product
	productInfor, err := app.DB.GetProductById(productId)
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	// lấy images liên quan tới product
	images, err := app.DB.GetInforImagesWithProductId(productId)
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	data := make(map[string]interface{})
	data["Categories"] = categories
	data["Brands"] = brands
	data["Tags"] = tags
	data["RelatedProducts"] = relatedProductsWithHostPort
	data["ProductInfo"] = productInfor
	data["Images"] = images
	// data["Quantity"] = cart.GetQuantity()

	if err := app.renderTemplate(w, r, "product-detail", &templateData{
		Data: data,
		Page: "PRODUCT",
	}, "right-column", "product"); err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}
}

// hàm render profile page
func (app *application) ProfilePage(w http.ResponseWriter, r *http.Request) {
	// lấy thông tin từ session
	session, err := app.Redis.Get(r, app.config.session_key)
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	// lấy email của người dùng
	var userEmail string
	if email, ok := session.Values["UserEmail"].(string); ok {
		userEmail = email
	} else {
		app.errorLog.Println("Không thể lấy thông tin người dùng")
		panic("Không thể lấy thông tin người dùng")
	}

	// lấy thông tin của người dùng
	user, err := app.DB.GetUserByEmail(userEmail)
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	// lấy danh sách sản phẩm đẵ đặt hàng
	boughtProducts, err := app.DB.GetBoughtProductWithUserId(user.ID)
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	// tính toán score
	count := 0
	for index := 0; index < len(boughtProducts); index++ {
		boughtProducts[index].No = count
		count++
	}

	// lấy danh sách shippingaddress của người dùng
	shippingAddress, err := app.DB.GetAddressForUser(user.ID)
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	// thêm số thứ tự cho address
	count = 1
	for index := 0; index < len(shippingAddress); index++ {
		shippingAddress[index].No = count
		count++
	}

	data := make(map[string]interface{})
	data["BoughtProducts"] = boughtProducts
	data["ShippingAddresses"] = shippingAddress
	data["Profile"] = user

	// render profile page
	if err := app.renderTemplate(w, r, "profile", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
		panic(err)
	}
}

// render contact page
func (app *application) ContactPage(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "contact", &templateData{
		Page: "CONTACT",
	}); err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}
}

// hàm xử lý yêu cầu dặt hàng sản phẩm
func (app *application) Order(w http.ResponseWriter, r *http.Request) {
	// body của reponse
	var payLoad struct {
		Message string `json:"message"`
	}

	// biến chứa thông tin từ request
	var data struct {
		AddressId string `json:"addressId"`
		Method    string `json:"method"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Mobile    string `json:"mobile"`
		Address   string `json:"address"`
		Country   string `json:"country"`
		City      string `json:"city"`
	}

	// đọc dữ liệu từ request
	err := app.readJSON(w, r, &data)
	if err != nil {
		app.errorLog.Println(err.Error())
		payLoad.Message = "Không thể đọc dữ liệu từ request"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// lấy session
	session, err := app.Redis.Get(r, app.config.session_key)
	if err != nil {
		app.errorLog.Println("Không lấy được session")
		payLoad.Message = "Không lấy được session"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// lấy email của người dùng trong session
	var userEmail string
	if email, ok := session.Values["UserEmail"].(string); ok {
		userEmail = email
	} else {
		app.errorLog.Println("Không tìm thấy thông tin người dùng")
		payLoad.Message = "Không tìm thấy thông tin người dùng"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// lấy thông tin của người dùng
	user, err := app.DB.GetUserByEmail(userEmail)
	if err != nil {
		app.errorLog.Println("Không tìm thấy thông tin người dùng")
		payLoad.Message = "Không tìm thấy thông tin người dùng"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// lấy address của người dùng
	var address models.Address
	if data.AddressId != "0" {
		addresses, err := app.DB.GetAddressForUser(user.ID)
		if err != nil {
			app.errorLog.Println("Không tìm thấy địa chỉ người dùng")
			payLoad.Message = "Không tìm thấy địa chỉ người dùng"
			app.writeJSON(w, 400, payLoad)
			return
		}

		for _, addr := range addresses {
			if addr.ID == data.AddressId {
				jsonData, _ := json.Marshal(addr)
				_ = json.Unmarshal(jsonData, &address)
				break
			}
		}
	} else {
		address.ID = "0"
		address.FirstName = data.FirstName
		address.LastName = data.LastName
		address.Email = data.Email
		address.Mobile = data.Mobile
		address.Address = data.Address
		address.Country = data.Country
		address.City = data.City
	}

	// lấy thông tin giỏ hàng của người dùng
	cartDetail, err := app.DB.GetCartForUser(user.ID)
	if err != nil {
		app.errorLog.Println("Không thể lấy thông tin giỏ hàng của người dùng")
		payLoad.Message = "Không thể lấy thông tin giỏ hàng của người dùng"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// fmt.Println(cartDetail)

	// xử lý truy vấn đặt hàng
	if data.AddressId == "0" { // nếu dặt ở địa chỉ khác
		_, err = app.DB.InsertOrder(user.ID, cartDetail, data.FirstName, data.LastName, data.Email, data.Mobile, data.Address, data.Country, data.City)
	} else { // nếu đặt ở địa chỉ có sẵn
		_, err = app.DB.InsertOrder(user.ID, cartDetail, address.FirstName, address.LastName, address.Email, address.Mobile, address.Address, address.Country, address.City)
	}

	// nếu thao tác truy vấn xảy ra lỗi
	if err != nil {
		app.errorLog.Println(err.Error())
		payLoad.Message = "Không thể xác nhận giao dịch"
		app.writeJSON(w, 400, payLoad)
		return
	}

	// dặt hàng thành công
	app.infoLog.Println("Xác nhận đặt hàng thành công")
	payLoad.Message = "Xác nhận đặt hàng thành công"
	app.writeJSON(w, 200, payLoad)
}
