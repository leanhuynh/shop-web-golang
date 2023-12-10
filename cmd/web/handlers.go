package main

import (
	"context"
	"database/sql"
	"fmt"
	"myapp/internal/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

func (app *application) HomePage(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// query for feature products
	rows, err := app.DB.DB.QueryContext(ctx, `
	select 
		id, name, image_path, stars, price, old_price, created_at
	from Product
	order by stars DESC
	limit ?
	`, 10)

	if err != nil {
		app.errorLog.Println(err)
	}

	var featureProducts []models.Product
	for rows.Next() {
		var product models.Product
		err = rows.Scan(&product.ID, &product.Name, &product.ImagePath, &product.Stars, &product.Price, &product.OldPrice, &product.CreatedAt)
		if err != nil {
			app.errorLog.Println(err)
		}
		featureProducts = append(featureProducts, product)
	}

	// query for category product
	rows, err = app.DB.DB.QueryContext(ctx, `
	select 
		id, name, image_path
	from Category
	`)
	if err != nil {
		app.errorLog.Println(err)
	}

	var CategoryArray []models.Category
	for rows.Next() {
		var category models.Category
		err = rows.Scan(&category.ID, &category.Name, &category.ImagePath)
		if err != nil {
			app.errorLog.Println(err)
		}
		CategoryArray = append(CategoryArray, category)
	}

	// query for recent products
	rows, err = app.DB.DB.QueryContext(ctx, `
	select 
		id, name, image_path, stars, price, old_price, created_at
	from Product
	order by created_at DESC
	limit ?
	`, 10)
	if err != nil {
		app.errorLog.Println(err)
	}

	var RecentProducts []models.Product
	for rows.Next() {
		var product models.Product
		err = rows.Scan(&product.ID, &product.Name, &product.ImagePath, &product.Stars, &product.Price, &product.OldPrice, &product.CreatedAt)
		if err != nil {
			app.errorLog.Println(err)
		}
		RecentProducts = append(RecentProducts, product)
	}

	//query for brands
	rows, err = app.DB.DB.QueryContext(ctx, `
	select 
		id, name, image_path
	from brand
	`)

	if err != nil {
		app.errorLog.Println(err)
	}

	var Brands []models.Brand
	for rows.Next() {
		var brand models.Brand
		err = rows.Scan(&brand.ID, &brand.Name, &brand.ImagePath)
		if err != nil {
			app.errorLog.Println(err)
		}
		Brands = append(Brands, brand)
	}

	data := make(map[string]interface{})
	data["featureProducts"] = featureProducts
	data["CategoryArray"] = CategoryArray
	data["RecentProducts"] = RecentProducts
	data["Brands"] = Brands

	if err := app.renderTemplate(w, r, "index", &templateData{
		Data: data,
	}, "product"); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) ShowCart(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cart := session.Get(r.Context(), "cart").(Cart)

	data := make(map[string]interface{})
	// query for cart's product
	// get the product's id
	for index, item := range cart.Items {
		var product models.Product

		row := app.DB.DB.QueryRowContext(ctx, `
		select 
			name, image_path
		from 
			product
		where
			id = ?
		`, item.ID)

		err := row.Scan(&product.Name, &product.ImagePath)
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		cart.Items[index].Name = product.Name
		cart.Items[index].ImagePath = product.ImagePath
	}

	fmt.Println(cart.GetCart())
	data["CartProducts"] = cart.GetCart()

	if err := app.renderTemplate(w, r, "cart", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) ShowProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// find Sort Type
	queryParams := r.URL.Query()
	var sort string
	if len(queryParams["sort"]) == 1 {
		sort = queryParams["sort"][0]
	}

	if check := ContainString([]string{"price", "newest", "popular"}, sort); !check {
		sort = "price"
	}

	var categoryId int
	if len(queryParams["category"]) == 1 {
		categoryId, _ = strconv.Atoi(queryParams["category"][0])
	}

	var brandId int
	if len(queryParams["brand"]) == 1 {
		brandId, _ = strconv.Atoi(queryParams["brand"][0])
	}

	stmt := "select id, name, image_path, stars, price, old_price from Product "

	if categoryId != 0 {
		if strings.Contains(stmt, "where") {
			stmt += "and product.category_id = ?"
		} else {
			stmt += "where product.category_id = ?"
		}
	}

	if brandId != 0 {
		if strings.Contains(stmt, "where") {
			stmt += "and product.brand_id = ?"
		} else {
			stmt += "where product.brand_id = ?"
		}
	}

	// query for List Of Product
	var rows *sql.Rows
	var err error

	if categoryId != 0 {
		rows, err = app.DB.DB.QueryContext(ctx, stmt, categoryId)
	} else if brandId != 0 {
		rows, err = app.DB.DB.QueryContext(ctx, stmt, brandId)
	} else {
		rows, err = app.DB.DB.QueryContext(ctx, stmt)
	}

	if err != nil {
		app.errorLog.Println(err)
	}

	var ProductList []models.Product
	for rows.Next() {
		var product models.Product
		err = rows.Scan(&product.ID, &product.Name, &product.ImagePath, &product.Stars, &product.Price, &product.OldPrice)
		if err != nil {
			app.errorLog.Println(err)
		}
		ProductList = append(ProductList, product)
	}

	// get original url
	OriginalUrl := r.URL.String()
	OriginalUrl = RemoveParams("sort", OriginalUrl)
	OriginalUrl += "?"

	// query for category
	rows, err = app.DB.DB.QueryContext(ctx, `
	select 
		t1.id, t1.name, t1.image_path, count(t2.id) as num_of_product  
	from 
		Category as t1
	left join 
		 product as t2 on t1.id = t2.category_id 
	group by 
		t1.id;
	`)

	if err != nil {
		app.errorLog.Println(err)
	}

	var Categories []models.Category

	for rows.Next() {
		var category models.Category
		err = rows.Scan(&category.ID, &category.Name, &category.ImagePath, &category.NumOfProduct)
		if err != nil {
			app.errorLog.Println(err)
		}
		Categories = append(Categories, category)
	}

	// query brands
	rows, err = app.DB.DB.QueryContext(ctx, `
	select 
		t1.id, t1.name, t1.image_path, count(t2.id) as num_of_product  
	from 
		Brand as t1
	left join 
		 Product as t2 
	on t1.id = t2.brand_id 
	group by 
		t1.id
	`)

	if err != nil {
		app.errorLog.Println(err)
	}

	var Brands []models.Brand

	for rows.Next() {
		var brand models.Brand
		err = rows.Scan(&brand.ID, &brand.Name, &brand.ImagePath, &brand.NumOfProduct)
		if err != nil {
			app.errorLog.Println(err)
		}
		Brands = append(Brands, brand)
	}

	// query for tags
	rows, err = app.DB.DB.QueryContext(ctx, `
	select 
		id, name
	from 
		Tag
	`)

	if err != nil {
		app.errorLog.Println(err)
	}

	var Tags []models.Tag

	for rows.Next() {
		var tag models.Tag
		err = rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			app.errorLog.Println(err)
		}
		Tags = append(Tags, tag)
	}

	data := make(map[string]interface{})
	data["ProductList"] = ProductList
	data["Sort"] = sort
	data["OriginalUrl"] = OriginalUrl
	data["Categories"] = Categories
	data["Brands"] = Brands
	data["Tags"] = Tags

	if err := app.renderTemplate(w, r, "product-list", &templateData{
		Data: data,
	}, "product", "right-column"); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) CheckOut(w http.ResponseWriter, r *http.Request) {
	// if err := app.renderTemplate(w, r, "cart", &templateData{}); err != nil {
	// 	app.errorLog.Println(err)
	// }
}

func (app *application) PlaceOrders(w http.ResponseWriter, r *http.Request) {

}

func (app *application) ShowProductDetail(w http.ResponseWriter, r *http.Request) {
	var productId int
	if chi.URLParam(r, "product_id") != "" {
		productId, _ = strconv.Atoi(chi.URLParam(r, "product_id"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// query for category
	rows, err := app.DB.DB.QueryContext(ctx, `
	select 
		t1.id, t1.name, t1.image_path, count(t2.id) as num_of_product  
	from 
		Category as t1
	left join 
		product as t2 on t1.id = t2.category_id 
	group by 
		t1.id;
	`)

	if err != nil {
		app.errorLog.Println(err)
	}

	var Categories []models.Category

	for rows.Next() {
		var category models.Category
		err = rows.Scan(&category.ID, &category.Name, &category.ImagePath, &category.NumOfProduct)
		if err != nil {
			app.errorLog.Println(err)
		}
		Categories = append(Categories, category)
	}

	// query brands
	rows, err = app.DB.DB.QueryContext(ctx, `
	select 
		t1.id, t1.name, t1.image_path, count(t2.id) as num_of_product  
	from 
		Brand as t1
	left join 
		Product as t2 
	on t1.id = t2.brand_id 
	group by 
		t1.id
	`)

	if err != nil {
		app.errorLog.Println(err)
	}

	var Brands []models.Brand

	for rows.Next() {
		var brand models.Brand
		err = rows.Scan(&brand.ID, &brand.Name, &brand.ImagePath, &brand.NumOfProduct)
		if err != nil {
			app.errorLog.Println(err)
		}
		Brands = append(Brands, brand)
	}

	// query for tags
	rows, err = app.DB.DB.QueryContext(ctx, `
	select 
		id, name
	from 
		Tag
	`)

	if err != nil {
		app.errorLog.Println(err)
	}

	var Tags []models.Tag

	for rows.Next() {
		var tag models.Tag
		err = rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			app.errorLog.Println(err)
		}
		Tags = append(Tags, tag)
	}

	// query for related products
	rows, err = app.DB.DB.QueryContext(ctx, `
	select
		t1.id, t1.name, t1.image_path, t1.old_price, t1.price
	from
		product as t1
	inner join
		producttag as t2
	where
		t1.id = t2.product_id and t1.id != ? and t2.tag_id in (
			select
				producttag.tag_id
			from
				producttag
			where
				producttag.product_id = ?
		)
	`, productId, productId)

	if err != nil {
		app.errorLog.Println(err)
	}

	var RelatedProducts []models.Product

	for rows.Next() {
		var product models.Product
		err = rows.Scan(&product.ID, &product.Name, &product.ImagePath, &product.OldPrice, &product.Price)
		if err != nil {
			app.errorLog.Println(err)
		}
		RelatedProducts = append(RelatedProducts, product)
	}

	// query infor product-detail
	row := app.DB.DB.QueryRowContext(ctx, `
	select 
		id, name, image_path, old_price, price, summary
	from Product
	where id = ?
	`, productId)

	if err != nil {
		app.errorLog.Println(err)
	}

	var ProductInfo models.Product
	err = row.Scan(
		&ProductInfo.ID,
		&ProductInfo.Name,
		&ProductInfo.ImagePath,
		&ProductInfo.OldPrice,
		&ProductInfo.Price,
		&ProductInfo.Summary,
	)

	if err != nil {
		app.errorLog.Println(err)
	}

	// query image-path
	rows, err = app.DB.DB.QueryContext(ctx, `
	select 
		Image.name, Image.image_path
	from Image
	where Image.product_id = ?
	`, productId)

	if err != nil {
		app.errorLog.Println(err)
	}

	var Images []models.Image

	for rows.Next() {
		var image models.Image
		err = rows.Scan(&image.Name, &image.ImagePath)
		if err != nil {
			app.errorLog.Println(err)
		}
		Images = append(Images, image)
	}

	data := make(map[string]interface{})
	data["Categories"] = Categories
	data["Brands"] = Brands
	data["Tags"] = Tags
	data["RelatedProducts"] = RelatedProducts
	data["ProductInfo"] = ProductInfo
	data["Images"] = Images

	if err := app.renderTemplate(w, r, "product-detail", &templateData{
		Data: data,
	}, "right-column", "product"); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) AddCart(w http.ResponseWriter, r *http.Request) {
	// get information of current cart
	cart := session.Get(r.Context(), "cart").(Cart)

	// get information of product wanted to be added
	var order struct {
		ProductId int `json:"id"`
		Quantity  int `json:"quantity"`
	}

	err := app.readJSON(w, r, &order)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	// handle to add product into cart

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := app.DB.DB.QueryRowContext(ctx, `
	select 
		price
	from 
		product
	where 
		id = ?
	`, order.ProductId)

	var price float32
	err = row.Scan(&price)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	productExisted := false
	for index, item := range cart.Items {
		// neu trong cart co product can them vao
		if item.ID == order.ProductId {
			cart.Items[index].Quantity += order.Quantity
			cart.Items[index].Price += price
			cart.Items[index].Total = float32(cart.Items[index].Quantity) * price
			productExisted = true
			break
		}
	}

	// neu trong cart khong co product can them vao
	if !productExisted {
		cart.Items = append(cart.Items, models.OrderDetail{
			ID:       order.ProductId,
			Quantity: order.Quantity,
			Price:    price,
			Total:    price * float32(order.Quantity),
		})
	}

	// send status back to client side
	session.Put(r.Context(), "cart", cart)
}

func (app *application) UpdateCart(w http.ResponseWriter, r *http.Request) {

}

func (app *application) RemoveProduct(w http.ResponseWriter, r *http.Request) {

}

func (app *application) ClearCart(w http.ResponseWriter, r *http.Request) {

}
