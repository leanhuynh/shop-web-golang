package main

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"myapp/internal/models"
// 	"myapp/internal/urlsigner"
// 	"net/http"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/go-chi/chi/v5"
// )

// func (app *application) HomePage(w http.ResponseWriter, r *http.Request) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	// query for feature products
// 	rows, err := app.DB.DB.QueryContext(ctx, `
// 	select
// 		id, name, image_path, stars, price, old_price, created_at
// 	from Product
// 	order by stars DESC
// 	limit ?
// 	`, 10)

// 	if err != nil {
// 		app.errorLog.Println(err)
// 	}

// 	var featureProducts []models.Product
// 	for rows.Next() {
// 		var product models.Product
// 		err = rows.Scan(&product.ID, &product.Name, &product.ImagePath, &product.Stars, &product.Price, &product.OldPrice, &product.CreatedAt)
// 		if err != nil {
// 			app.errorLog.Println(err)
// 		}
// 		featureProducts = append(featureProducts, product)
// 	}

// 	// query for category product
// 	rows, err = app.DB.DB.QueryContext(ctx, `
// 	select
// 		id, name, image_path
// 	from Category
// 	`)
// 	if err != nil {
// 		app.errorLog.Println(err)
// 	}

// 	var CategoryArray []models.Category
// 	for rows.Next() {
// 		var category models.Category
// 		err = rows.Scan(&category.ID, &category.Name, &category.ImagePath)
// 		if err != nil {
// 			app.errorLog.Println(err)
// 		}
// 		CategoryArray = append(CategoryArray, category)
// 	}

// 	// query for recent products
// 	rows, err = app.DB.DB.QueryContext(ctx, `
// 	select
// 		id, name, image_path, stars, price, old_price, created_at
// 	from Product
// 	order by created_at DESC
// 	limit ?
// 	`, 10)
// 	if err != nil {
// 		app.errorLog.Println(err)
// 	}

// 	var RecentProducts []models.Product
// 	for rows.Next() {
// 		var product models.Product
// 		err = rows.Scan(&product.ID, &product.Name, &product.ImagePath, &product.Stars, &product.Price, &product.OldPrice, &product.CreatedAt)
// 		if err != nil {
// 			app.errorLog.Println(err)
// 		}
// 		RecentProducts = append(RecentProducts, product)
// 	}

// 	//query for brands
// 	rows, err = app.DB.DB.QueryContext(ctx, `
// 	select
// 		id, name, image_path
// 	from brand
// 	`)

// 	if err != nil {
// 		app.errorLog.Println(err)
// 	}

// 	var Brands []models.Brand
// 	for rows.Next() {
// 		var brand models.Brand
// 		err = rows.Scan(&brand.ID, &brand.Name, &brand.ImagePath)
// 		if err != nil {
// 			app.errorLog.Println(err)
// 		}
// 		Brands = append(Brands, brand)
// 	}

// 	// query for quantity of cart
// 	// get information of current cart
// 	// cart := session.Get(r.Context(), "cart").(Cart)

// 	data := make(map[string]interface{})
// 	data["featureProducts"] = featureProducts
// 	data["CategoryArray"] = CategoryArray
// 	data["RecentProducts"] = RecentProducts
// 	data["Brands"] = Brands
// 	// data["Quantity"] = cart.GetQuantity()

// 	if err := app.renderTemplate(w, r, "index", &templateData{
// 		Data: data,
// 	}, "product"); err != nil {
// 		app.errorLog.Println(err)
// 	}
// }

// func (app *application) Register(w http.ResponseWriter, r *http.Request) {
// 	if err := app.renderTemplate(w, r, "register", &templateData{}); err != nil {
// 		app.errorLog.Println(err)
// 	}
// }

// func (app *application) Login(w http.ResponseWriter, r *http.Request) {
// 	if err := app.renderTemplate(w, r, "login", &templateData{}); err != nil {
// 		app.errorLog.Println(err)
// 	}
// }

// func (app *application) ForgotPassword(w http.ResponseWriter, r *http.Request) {
// 	if err := app.renderTemplate(w, r, "forgot-password", &templateData{}); err != nil {
// 		app.errorLog.Println(err)
// 	}
// }

// func (app *application) ResetPassword(w http.ResponseWriter, r *http.Request) {
// 	email := r.URL.Query().Get("email")
// 	theURL := r.RequestURI
// 	testURL := fmt.Sprintf("%s%s", "http://localhost:4000", theURL)

// 	signer := urlsigner.Signer{
// 		Secret: []byte("secret"),
// 	}

// 	valid := signer.VerifyToken(testURL)

// 	if !valid {
// 		app.errorLog.Println("Invalid url - tampering detected")
// 		return
// 	}

// 	// make sure not expired
// 	expired := signer.Expired(testURL, 60)
// 	if expired {
// 		app.errorLog.Println("Link expired")
// 		return
// 	}

// 	data := make(map[string]interface{})
// 	data["email"] = email

// 	if err := app.renderTemplate(w, r, "reset-password", &templateData{}); err != nil {
// 		app.errorLog.Println(err)
// 	}
// }

// func (app *application) ShowCart(w http.ResponseWriter, r *http.Request) {
// 	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	// defer cancel()

// 	// cart := session.Get(r.Context(), "cart").(Cart)

// 	// data := make(map[string]interface{})
// 	// // query for cart's product
// 	// // get the product's id
// 	// for index, item := range cart.Items {
// 	// 	var product models.Product

// 	// 	row := app.DB.DB.QueryRowContext(ctx, `
// 	// 	select
// 	// 		name, image_path
// 	// 	from
// 	// 		product
// 	// 	where
// 	// 		id = ?
// 	// 	`, item.ID)

// 	// 	err := row.Scan(&product.Name, &product.ImagePath)
// 	// 	if err != nil {
// 	// 		app.errorLog.Println(err)
// 	// 		return
// 	// 	}
// 	// 	cart.Items[index].Name = product.Name
// 	// 	cart.Items[index].ImagePath = product.ImagePath
// 	// }

// 	// data["Quantity"] = cart.GetQuantity()
// 	// data["CartProducts"] = cart.GetCart()

// 	if err := app.renderTemplate(w, r, "cart", &templateData{
// 		// Data: data,
// 	}); err != nil {
// 		app.errorLog.Println(err)
// 	}
// }

// func (app *application) ShowProduct(w http.ResponseWriter, r *http.Request) {
// 	stmt := "select t1.id, t1.name, t1.image_path, t1.stars, t1.price, t1.old_price from Product t1 "

// 	// find Sort Type
// 	queryParams := r.URL.Query()

// 	// find category type
// 	var categoryId int
// 	if len(queryParams["category"]) == 1 {
// 		categoryId, _ = strconv.Atoi(queryParams["category"][0])
// 	}
// 	if categoryId < 0 {
// 		categoryId = 0
// 	}

// 	// find brand type
// 	var brandId int
// 	if len(queryParams["brand"]) == 1 {
// 		brandId, _ = strconv.Atoi(queryParams["brand"][0])
// 	}
// 	if brandId < 0 {
// 		brandId = 0
// 	}

// 	var sort string
// 	if len(queryParams["sort"]) == 1 {
// 		sort = queryParams["sort"][0]
// 	}

// 	if check := ContainString([]string{"price", "newest", "popular"}, sort); !check {
// 		sort = "price"
// 	}

// 	// find current page
// 	limit := 6
// 	var page int
// 	if len(queryParams["page"]) == 1 {
// 		page, _ = strconv.Atoi(queryParams["page"][0])
// 	}
// 	fmt.Println(page)
// 	fmt.Println(categoryId)
// 	// if not page params in url
// 	if page < 1 {
// 		page = 1
// 	}

// 	// find search keyword
// 	var search string
// 	if len(queryParams["search"]) == 1 {
// 		search = queryParams["search"][0]
// 	}

// 	// find tag_id
// 	var tagId int
// 	if len(queryParams["tag"]) == 1 {
// 		tagId, _ = strconv.Atoi(queryParams["tag"][0])
// 	}
// 	if tagId < 0 {
// 		tagId = 0
// 	}

// 	// add tag_id to query
// 	// just check if tag_id is exist
// 	// if not, not add query
// 	if tagId != 0 {
// 		stmt += "inner join producttag t2 on t1.id = t2.product_id where t2.tag_id = ? "
// 	}

// 	// add keyword (search) to query
// 	if search != "" {
// 		if strings.Contains(stmt, "where") {
// 			stmt += "and t1.name like ? "
// 		} else {
// 			stmt += "where t1.name like ? "
// 		}
// 	}

// 	// add category_id to query
// 	if categoryId != 0 {
// 		if strings.Contains(stmt, "where") {
// 			stmt += "and t1.category_id = ? "
// 		} else {
// 			stmt += "where t1.category_id = ? "
// 		}
// 	}

// 	// add brand_id to query
// 	if brandId != 0 {
// 		if strings.Contains(stmt, "where") {
// 			stmt += "and t1.brand_id = ? "
// 		} else {
// 			stmt += "where t1.brand_id = ? "
// 		}
// 	}

// 	// add sort to query
// 	stmt += "order by "
// 	if sort == "price" {
// 		stmt += "t1.price DESC "
// 	} else if sort == "newest" {
// 		stmt += "t1.created_at DESC "
// 	} else if sort == "popular" {
// 		stmt += "t1.stars DESC "
// 	}

// 	// query for List Of Product
// 	queryInfor := []interface{}{}
// 	if tagId != 0 {
// 		queryInfor = append(queryInfor, tagId)
// 	}

// 	if search != "" {
// 		queryInfor = append(queryInfor, "%"+search+"%")
// 	}

// 	if categoryId != 0 {
// 		queryInfor = append(queryInfor, categoryId)
// 	}

// 	if brandId != 0 {
// 		queryInfor = append(queryInfor, brandId)
// 	}

// 	// query de biet total row
// 	ProductList, err := app.DB.GetProductList(stmt, queryInfor)
// 	if err != nil {
// 		app.errorLog.Println(err)
// 		return
// 	}
// 	PageSize := (len(ProductList) + limit - 1) / limit

// 	// add page and limit
// 	stmt += "limit ? offset ?"

// 	// them vao limit va offset = (currentpage - 1) * limit
// 	queryInfor = append(queryInfor, limit)
// 	queryInfor = append(queryInfor, (page-1)*limit)

// 	ProductList, err = app.DB.GetProductList(stmt, queryInfor)

// 	if err != nil {
// 		app.errorLog.Println(err)
// 		return
// 	}

// 	// get original url
// 	OriginalUrl := r.URL.String()
// 	OriginalUrl = RemoveParams("sort", OriginalUrl)
// 	OriginalUrl += "?"

// 	// query for category
// 	Categories, err := app.DB.GetAllCategories()
// 	if err != nil {
// 		app.errorLog.Println(err)
// 		return
// 	}

// 	// query brands
// 	Brands, err := app.DB.GetAllBrand()
// 	if err != nil {
// 		app.errorLog.Println(err)
// 		return
// 	}

// 	// query for tags
// 	Tags, err := app.DB.GetAllTag()
// 	if err != nil {
// 		app.errorLog.Println(err)
// 		return
// 	}

// 	// query for quantity of cart
// 	// get information of current cart
// 	// cart := session.Get(r.Context(), "cart").(Cart)

// 	data := make(map[string]interface{})
// 	data["ProductList"] = ProductList
// 	data["Sort"] = sort
// 	data["OriginalUrl"] = OriginalUrl
// 	data["Categories"] = Categories
// 	data["Brands"] = Brands
// 	data["Tags"] = Tags
// 	// data["Quantity"] = cart.GetQuantity()
// 	data["CurrentPage"] = page
// 	data["PageSize"] = PageSize

// 	if err := app.renderTemplate(w, r, "product-list", &templateData{
// 		Data: data,
// 	}, "product", "right-column"); err != nil {
// 		app.errorLog.Println(err)
// 	}
// }

// func (app *application) CheckOut(w http.ResponseWriter, r *http.Request) {
// 	if err := app.renderTemplate(w, r, "checkout", &templateData{}); err != nil {
// 		app.errorLog.Println(err)
// 	}
// }

// // func (app *application) ChargeCard(w http.ResponseWriter, r *http.Request) {
// // 	if err := app.renderTemplate(w, r, "terminal", &templateData{}, "stripe-js"); err != nil {
// // 		app.errorLog.Println(err)
// // 	}
// // }

// func (app *application) ShowProductDetail(w http.ResponseWriter, r *http.Request) {
// 	var productId int
// 	if chi.URLParam(r, "product_id") != "" {
// 		productId, _ = strconv.Atoi(chi.URLParam(r, "product_id"))
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	// query for category
// 	rows, err := app.DB.DB.QueryContext(ctx, `
// 	select
// 		t1.id, t1.name, t1.image_path, count(t2.id) as num_of_product
// 	from
// 		Category as t1
// 	left join
// 		product as t2 on t1.id = t2.category_id
// 	group by
// 		t1.id;
// 	`)

// 	if err != nil {
// 		app.errorLog.Println(err)
// 	}

// 	var Categories []models.Category

// 	for rows.Next() {
// 		var category models.Category
// 		err = rows.Scan(&category.ID, &category.Name, &category.ImagePath, &category.NumOfProduct)
// 		if err != nil {
// 			app.errorLog.Println(err)
// 		}
// 		Categories = append(Categories, category)
// 	}

// 	// query brands
// 	rows, err = app.DB.DB.QueryContext(ctx, `
// 	select
// 		t1.id, t1.name, t1.image_path, count(t2.id) as num_of_product
// 	from
// 		Brand as t1
// 	left join
// 		Product as t2
// 	on t1.id = t2.brand_id
// 	group by
// 		t1.id
// 	`)

// 	if err != nil {
// 		app.errorLog.Println(err)
// 	}

// 	var Brands []models.Brand

// 	for rows.Next() {
// 		var brand models.Brand
// 		err = rows.Scan(&brand.ID, &brand.Name, &brand.ImagePath, &brand.NumOfProduct)
// 		if err != nil {
// 			app.errorLog.Println(err)
// 		}
// 		Brands = append(Brands, brand)
// 	}

// 	// query for tags
// 	rows, err = app.DB.DB.QueryContext(ctx, `
// 	select
// 		id, name
// 	from
// 		Tag
// 	`)

// 	if err != nil {
// 		app.errorLog.Println(err)
// 	}

// 	var Tags []models.Tag

// 	for rows.Next() {
// 		var tag models.Tag
// 		err = rows.Scan(&tag.ID, &tag.Name)
// 		if err != nil {
// 			app.errorLog.Println(err)
// 		}
// 		Tags = append(Tags, tag)
// 	}

// 	// query for related products
// 	rows, err = app.DB.DB.QueryContext(ctx, `
// 	select
// 		t1.id, t1.name, t1.image_path, t1.old_price, t1.price
// 	from
// 		product as t1
// 	inner join
// 		producttag as t2
// 	where
// 		t1.id = t2.product_id and t1.id != ? and t2.tag_id in (
// 			select
// 				producttag.tag_id
// 			from
// 				producttag
// 			where
// 				producttag.product_id = ?
// 		)
// 	`, productId, productId)

// 	if err != nil {
// 		app.errorLog.Println(err)
// 	}

// 	var RelatedProducts []models.Product

// 	for rows.Next() {
// 		var product models.Product
// 		err = rows.Scan(&product.ID, &product.Name, &product.ImagePath, &product.OldPrice, &product.Price)
// 		if err != nil {
// 			app.errorLog.Println(err)
// 		}
// 		RelatedProducts = append(RelatedProducts, product)
// 	}

// 	// query infor product-detail
// 	row := app.DB.DB.QueryRowContext(ctx, `
// 	select
// 		id, name, image_path, old_price, price, summary
// 	from Product
// 	where id = ?
// 	`, productId)

// 	if err != nil {
// 		app.errorLog.Println(err)
// 	}

// 	var ProductInfo models.Product
// 	err = row.Scan(
// 		&ProductInfo.ID,
// 		&ProductInfo.Name,
// 		&ProductInfo.ImagePath,
// 		&ProductInfo.OldPrice,
// 		&ProductInfo.Price,
// 		&ProductInfo.Summary,
// 	)

// 	if err != nil {
// 		app.errorLog.Println(err)
// 	}

// 	// query image-path
// 	rows, err = app.DB.DB.QueryContext(ctx, `
// 	select
// 		Image.name, Image.image_path
// 	from Image
// 	where Image.product_id = ?
// 	`, productId)

// 	if err != nil {
// 		app.errorLog.Println(err)
// 	}

// 	var Images []models.Image

// 	for rows.Next() {
// 		var image models.Image
// 		err = rows.Scan(&image.Name, &image.ImagePath)
// 		if err != nil {
// 			app.errorLog.Println(err)
// 		}
// 		Images = append(Images, image)
// 	}

// 	// query for quantity of cart
// 	// get information of current cart
// 	// cart := session.Get(r.Context(), "cart").(Cart)

// 	data := make(map[string]interface{})
// 	data["Categories"] = Categories
// 	data["Brands"] = Brands
// 	data["Tags"] = Tags
// 	data["RelatedProducts"] = RelatedProducts
// 	data["ProductInfo"] = ProductInfo
// 	data["Images"] = Images
// 	// data["Quantity"] = cart.GetQuantity()

// 	if err := app.renderTemplate(w, r, "product-detail", &templateData{
// 		Data: data,
// 	}, "right-column", "product"); err != nil {
// 		app.errorLog.Println(err)
// 	}
// }

// // func (app *application) AddCart(w http.ResponseWriter, r *http.Request) {
// // 	// get information of current cart
// // 	cart := session.Get(r.Context(), "cart").(Cart)

// // 	// get information of product wanted to be added
// // 	var infor struct {
// // 		ProductId int `json:"id"`
// // 		Quantity  int `json:"quantity"`
// // 	}

// // 	var payLoad struct {
// // 		Message string `json:"message"`
// // 	}

// // 	err := app.readJSON(w, r, &infor)
// // 	if err != nil {
// // 		payLoad.Message = "Can not read information"
// // 		app.writeJSON(w, 400, payLoad)
// // 		app.errorLog.Println(err)
// // 		return
// // 	}

// // 	// query product's information
// // 	product, err := app.DB.GetProductById(infor.ProductId)
// // 	if err != nil {
// // 		payLoad.Message = "Product's id not found"
// // 		app.writeJSON(w, 400, payLoad)
// // 		app.errorLog.Println(err)
// // 		return
// // 	}

// // 	// check bonous quantity is not larger then product's quantity
// // 	if infor.Quantity > product.Quantity {
// // 		payLoad.Message = "Vuot qua so luong hien co"
// // 		app.writeJSON(w, 400, payLoad)
// // 		app.errorLog.Println(err)
// // 		return
// // 	}

// // 	// them vao cart
// // 	cart.AddCart(product, infor.Quantity)

// // 	// send status back to client side
// // 	session.Put(r.Context(), "cart", cart)

// // 	payLoad.Message = "Add product successfully"
// // 	app.writeJSON(w, 200, payLoad)
// // }

// // func (app *application) UpdateCart(w http.ResponseWriter, r *http.Request) {
// // 	// get information of current cart
// // 	cart := session.Get(r.Context(), "cart").(Cart)

// // 	// get information of product wanted to be added
// // 	var infor struct {
// // 		ProductId int `json:"id"`
// // 		Quantity  int `json:"quantity"`
// // 	}

// // 	var payLoad struct {
// // 		Message string `json:"message"`
// // 	}

// // 	err := app.readJSON(w, r, &infor)
// // 	if err != nil {
// // 		payLoad.Message = "Can not read information"
// // 		app.writeJSON(w, 400, payLoad)
// // 		app.errorLog.Println(err)
// // 		return
// // 	}

// // 	// query product's information
// // 	product, err := app.DB.GetProductById(infor.ProductId)
// // 	if err != nil {
// // 		payLoad.Message = "Product's id not found"
// // 		app.writeJSON(w, 400, payLoad)
// // 		app.errorLog.Println(err)
// // 		return
// // 	}

// // 	// check quantity khong lon hon so luong san pham hien co
// // 	if infor.Quantity > product.Quantity {
// // 		payLoad.Message = "Vuot qua so luong san pham hien co"
// // 		app.writeJSON(w, 400, payLoad)
// // 		app.errorLog.Println(err)
// // 		return
// // 	}

// // 	// update cart
// // 	cart.UpdateCart(product, infor.Quantity)

// // 	// cap nhat trong session
// // 	session.Put(r.Context(), "cart", cart.GetCart())

// // 	// send to client
// // 	payLoad.Message = "Update cart successfully"
// // 	app.writeJSON(w, 200, payLoad)
// // }

// // func (app *application) RemoveCart(w http.ResponseWriter, r *http.Request) {
// // 	// get information of current cart
// // 	cart := session.Get(r.Context(), "cart").(Cart)

// // 	// get information of product wanted to be added
// // 	var infor struct {
// // 		ProductId int `json:"id"`
// // 		Quantity  int `json:"quantity"`
// // 	}

// // 	var payLoad struct {
// // 		Message string `json:"message"`
// // 	}

// // 	err := app.readJSON(w, r, &infor)
// // 	if err != nil {
// // 		payLoad.Message = "Can not read information"
// // 		app.writeJSON(w, 400, payLoad)
// // 		app.errorLog.Println(err)
// // 		return
// // 	}

// // 	// remove cart
// // 	cart.RemoveCart(infor.ProductId)

// // 	// update session
// // 	session.Put(r.Context(), "cart", cart)

// // 	// send response
// // 	payLoad.Message = "Remove product successfully"
// // 	app.writeJSON(w, 200, payLoad)
// // }

// // func (app *application) ClearCart(w http.ResponseWriter, r *http.Request) {
// // 	// get information of product wanted to be added
// // 	var infor struct {
// // 		ProductId int `json:"id"`
// // 		Quantity  int `json:"quantity"`
// // 	}

// // 	var payLoad struct {
// // 		Message string `json:"message"`
// // 	}

// // 	err := app.readJSON(w, r, &infor)
// // 	if err != nil {
// // 		payLoad.Message = "Can not read information"
// // 		app.writeJSON(w, 400, payLoad)
// // 		app.errorLog.Println(err)
// // 		return
// // 	}

// // 	// clear cart
// // 	session.Put(r.Context(), "cart", &Cart{Quantity: 0})

// // 	// send response
// // 	payLoad.Message = "Remove product successfully"
// // 	app.writeJSON(w, 200, payLoad)
// // }

// // func (app *application) UpdateCart(w http.ResponseWriter, r *http.Request) {
// // 	// get information of current cart
// // 	cart := session.Get(r.Context(), "cart").(Cart)

// // 	// get information of product wanted to be update
// // 	var order struct {
// // 		ProductId int `json:"id"`
// // 		Quantity  int `json:"quantity"`
// // 	}

// // 	err := app.readJSON(w, r, &order)
// // 	if err != nil {
// // 		app.errorLog.Println(err)
// // 		return
// // 	}
// // 	// handle to add product into cart

// // 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// // 	defer cancel()

// // 	row := app.DB.DB.QueryRowContext(ctx, `
// // 	select
// // 		price
// // 	from
// // 		product
// // 	where
// // 		id = ?
// // 	`, order.ProductId)

// // 	var price float32
// // 	err = row.Scan(&price)
// // 	if err != nil {
// // 		app.errorLog.Println(err)
// // 		return
// // 	}

// // 	for index, item := range cart.Items {
// // 		// neu trong cart co product can them vao
// // 		if item.ID == order.ProductId {
// // 			if order.Quantity == 0 {
// // 				cart.deleteItemWithIndex(index)
// // 			} else if order.Quantity > 0 {
// // 				cart.Items[index].Quantity = order.Quantity
// // 				cart.Items[index].Price = price
// // 				cart.Items[index].Total = float32(cart.Items[index].Quantity) * price
// // 				break
// // 			}
// // 		}
// // 	}

// // 	// send status back to client side
// // 	session.Put(r.Context(), "cart", cart)
// // }

// // func (app *application) RemoveCart(w http.ResponseWriter, r *http.Request) {
// // 	// get information of current cart
// // 	cart := session.Get(r.Context(), "cart").(Cart)

// // 	// get information of product wanted to be update
// // 	var order struct {
// // 		ProductId int `json:"id"`
// // 	}

// // 	err := app.readJSON(w, r, &order)
// // 	if err != nil {
// // 		app.errorLog.Println(err)
// // 		return
// // 	}

// // 	for index, item := range cart.Items {
// // 		if item.ID == order.ProductId {
// // 			cart.deleteItemWithIndex(index)
// // 		}
// // 	}

// // 	session.Put(r.Context(), "cart", cart)
// // }

// // func (app *application) ClearCart(w http.ResponseWriter, r *http.Request) {
// // 	session.Put(r.Context(), "cart", &Cart{Quantity: 0})
// // }

// func (app *application) authenticateToken(r *http.Request) (*models.User, error) {
// 	authorizationHeader := r.Header.Get("Authorization")
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

// // ham tao trang profile
// func (app *application) Profile(w http.ResponseWriter, r *http.Request) {
// 	if err := app.renderTemplate(w, r, "profile", &templateData{}); err != nil {
// 		app.errorLog.Println(err)
// 	}
// }

// // render contact page
// func (app *application) Contact(w http.ResponseWriter, r *http.Request) {
// 	if err := app.renderTemplate(w, r, "contact", &templateData{}); err != nil {
// 		app.errorLog.Println(err)
// 	}
// }

// func (app *application) ErrorPage(w http.ResponseWriter, r *http.Request) {
// 	if err := app.renderTemplate(w, r, "error", &templateData{}); err != nil {
// 		app.errorLog.Println(err)
// 	}
// }
