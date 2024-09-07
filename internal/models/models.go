package models

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type DBModel struct {
	DB *sql.DB
}

type Models struct {
	DB DBModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Mobile    string `json:"mobile"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	IsAdmin   bool   `json:"is_admin"`
}

type Address struct {
	No        int    `json:"no"`
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Mobile    string `json:"mobile"`
	Address   string `json:"address"`
	Country   string `json:"country"`
	City      string `json:"city"`
	State     string `json:"state"`
	ZipCode   string `json:"zip_code"`
	IsDefault bool   `json:"is_default"`
}

type Brand struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ImagePath    string `json:"image_path"`
	NumOfProduct int    `json:"num_of_product"`
}

type Category struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ImagePath    string `json:"image_path"`
	NumOfProduct int    `json:"num_of_product"`
}

type Image struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ImagePath string `json:"image_path"`
	ProductId string `json:"product_id"`
}

type Message struct {
	ID      int    `json:"id"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type Order struct {
	ID       int     `json:"id"`
	Quantity int     `json:"quantity"`
	Total    float32 `json:"total"`
	SubTotal float32 `json:"sub_total"`
	Shipping float32 `json:"shipping"`
}

// type OrderDetail struct {
// 	ID        int     `json:"id"`
// 	Name      string  `json:"name"`
// 	ImagePath string  `json:"image_path"`
// 	Quantity  int     `json:"quantity"`
// 	Price     float32 `json:"price"`
// 	Total     float32 `json:"total"`
// }

type PlaceOrder struct {
	ID       int     `json:"id"`
	UserId   int     `json:"user_id"`
	Quantity int     `json:"quantity"`
	Price    float32 `json:"price"`
}

type Product struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	ImagePath     string    `json:"image_path"`
	OldPrice      float32   `json:"old_price"`
	Price         float32   `json:"price"`
	Summary       string    `json:"summary"`
	Description   string    `json:"description"`
	Specification string    `json:"specification"`
	Stars         float32   `json:"stars"`
	Quantity      int       `json:"quantity"`
	CategoryId    int       `json:"category_id"`
	BrandId       int       `json:"brand_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type ProductWithHostPort struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	ImagePath     string    `json:"image_path"`
	OldPrice      float32   `json:"old_price"`
	Price         float32   `json:"price"`
	Summary       string    `json:"summary"`
	Description   string    `json:"description"`
	Specification string    `json:"specification"`
	Stars         float32   `json:"stars"`
	Quantity      int       `json:"quantity"`
	CategoryId    int       `json:"category_id"`
	BrandId       int       `json:"brand_id"`
	CreatedAt     time.Time `json:"created_at"`
	Host          string    `json:"host"`
	Port          int       `json:"port"`
}

type ProductTag struct {
	TagId     int `json:"tag_id"`
	ProductId int `json:"product_id"`
}

type Review struct {
	ID     int     `json:"id"`
	Review string  `json:"review"`
	Stars  float32 `json:"stars"`
}

type Subscribe struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type WishList struct {
	ID        int `json:"id"`
	UserId    int `json:"user_id"`
	ProductId int `json:"product_id"`
}

type Invoice struct {
	No          int       `json:"no"`
	UserId      int       `json:"user_id"`
	ProductId   int       `json:"product_id"`
	ProductName string    `json:"product_name"`
	ImagePath   string    `json:"image_path"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
	Status      string    `json:"status"`
	BoughtDate  time.Time `json:"bought_date"`
}

func (m *DBModel) CreateUser(email string, password string, firstname string, lastname string, mobile string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// kiểm tra email có trùng hay không
	User, _ := m.GetUserByEmail(email)
	if User.Email == email {
		return errors.New("tồn tại tài khoản cùng email")
	}

	stmt := `insert into user (email, password, first_name, last_name, mobile) values (?, ?, ?, ?, ?)`
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = m.DB.ExecContext(ctx, stmt, email, hashedPassword, firstname, lastname, mobile)

	return err
}

func (m *DBModel) GetUserByEmail(email string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User

	stmt := `select id, email, password, first_name, last_name, is_admin, mobile from user where email = ?`
	row := m.DB.QueryRowContext(ctx, stmt, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.IsAdmin,
		&user.Mobile,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (m *DBModel) UpdatePasswordForUser(email string, new_password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `update user set user.password = ? where user.email = ?`

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(new_password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = m.DB.ExecContext(ctx, stmt, hashedPassword, email)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) GetUserById(id int) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// query to user with given id
	stmt := `select id, first_name, last_name, mobile, email, password, is_admin from user where user.id = ?`
	row := m.DB.QueryRowContext(ctx, stmt, id)

	// json
	var user User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.IsAdmin,
		&user.Mobile,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (m *DBModel) DeleteUserForToken(token string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// thuc hien query
	tokenHash := sha256.Sum256([]byte(token))
	stmt := `delete from token where token_hash = ?`
	_, err := m.DB.ExecContext(ctx, stmt, tokenHash)

	// neu khong delete thanh cong
	if err != nil {
		return false, err
	}

	// delete thanh cong
	return true, err
}

func (m *DBModel) GetAddressForUser(user_id int) ([]Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var Addresses []Address
	stmt := `select id, first_name, last_name, email, mobile, address, country, city, state, zip_code, is_default from address where user_id = ?`
	rows, err := m.DB.QueryContext(ctx, stmt, user_id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var address Address
		err := rows.Scan(
			&address.ID,
			&address.FirstName,
			&address.LastName,
			&address.Email,
			&address.Mobile,
			&address.Address,
			&address.Country,
			&address.City,
			&address.State,
			&address.ZipCode,
			&address.IsDefault,
		)
		if err != nil {
			return nil, err
		}
		Addresses = append(Addresses, address)
	}
	return Addresses, nil
}

// hàm xử lý order, thêm sản phẩm vào order
func (m *DBModel) InsertOrder(userId int, cart []CartDetail, firstName string, lastName string, email string, mobile string, address string, country string, city string) (bool, error) {
	if len(cart) == 0 {
		return false, errors.New("không có sản phẩm trong giỏ hàng")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	queryInfor := []interface{}{}
	stmt := "insert into Invoice (user_id, product_id, image_path, price, quantity, first_name, last_name, email, mobile, address, country, city) values "

	// thực hiện với product-0
	queryInfor = append(queryInfor, userId)
	queryInfor = append(queryInfor, cart[0].ProductId)
	queryInfor = append(queryInfor, cart[0].ImagePath)
	queryInfor = append(queryInfor, cart[0].Price)
	queryInfor = append(queryInfor, cart[0].Quantity)
	queryInfor = append(queryInfor, firstName)
	queryInfor = append(queryInfor, lastName)
	queryInfor = append(queryInfor, email)
	queryInfor = append(queryInfor, mobile)
	queryInfor = append(queryInfor, address)
	queryInfor = append(queryInfor, country)
	queryInfor = append(queryInfor, city)
	stmt += " (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) "

	for index := 1; index < len(cart); index++ {
		queryInfor = append(queryInfor, userId)
		queryInfor = append(queryInfor, cart[index].ProductId)
		queryInfor = append(queryInfor, cart[index].ImagePath)
		queryInfor = append(queryInfor, cart[index].Price)
		queryInfor = append(queryInfor, cart[index].Quantity)
		queryInfor = append(queryInfor, firstName)
		queryInfor = append(queryInfor, lastName)
		queryInfor = append(queryInfor, email)
		queryInfor = append(queryInfor, mobile)
		queryInfor = append(queryInfor, address)
		queryInfor = append(queryInfor, country)
		queryInfor = append(queryInfor, city)
		stmt += " , (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) "
	}

	fmt.Println(queryInfor...)
	fmt.Println(stmt)

	_, err := m.DB.ExecContext(ctx, stmt, queryInfor...)
	if err != nil { // nếu thao tác không thể thực thi
		return false, err
	}

	// thao tác được thực hiện --> xóa tất cả sản phẩm trong giỏ hàng
	err = m.RemoveCart(userId)
	if err != nil {
		return false, err
	}

	// thao tác được thực hiện --> cập nhật số lượng hiện có
	err = m.UpdateQuantityofProducts(cart)
	if err != nil {
		return false, err
	}

	return true, nil
}

// hàm cập nhật số lượng sản phẩm trong giỏ hàng theo danh sách
func (m *DBModel) UpdateQuantityofProducts(cart []CartDetail) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for index := 0; index < len(cart); index++ {
		// thực thi câu truy vấn cập nhật số lượng sản phẩm trong kho
		_, err := m.DB.ExecContext(ctx, `
			update product
			set quantity = quantity - ?
			where id = ?
		`, cart[index].Quantity, cart[index].ProductId)

		// nếu có lỗi xảy ra trong quá trình thực thi
		if err != nil {
			return err
		}
	}

	// nếu không có lỗi xảy ra
	return nil
}

func (m *DBModel) GetAllTags() ([]Tag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, `
	select 
		id, name
	from 
		Tag
	`)

	if err != nil {
		return nil, err
	}

	var Tags []Tag

	for rows.Next() {
		var tag Tag
		err = rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return nil, err
		}
		Tags = append(Tags, tag)
	}

	return Tags, nil
}

func (m *DBModel) GetProductById(product_id int) (Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// thuc hien query
	var product Product
	row := m.DB.QueryRowContext(ctx, `
	select 
		id, name, image_path, old_price, price, summary, description, specification, stars, quantity, brand_id, category_id
	from 
		product
	where
		id = ?
	`, product_id)

	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.ImagePath,
		&product.OldPrice,
		&product.Price,
		&product.Summary,
		&product.Description,
		&product.Specification,
		&product.Stars,
		&product.Quantity,
		&product.BrandId,
		&product.CategoryId,
	)

	// kiem tra ket qua
	if err != nil {
		return product, err
	}

	return product, nil
}

func (m *DBModel) GetAvailableShippingAddress(user User) ([]Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// thực hiện query lấy tất cả Available Shipping Address của User
	rows, err := m.DB.QueryContext(ctx, `
	select user_id, first_name, last_name, email, mobile, address, country, city, state, zip_code, is_default 
	from Address
	where user_id = ?
	`, user.ID)

	// kiểm tra nếu có lỗi xảy ra
	if err != nil {
		return nil, err
	}

	var Addresses []Address

	for rows.Next() {
		var address Address
		err = rows.Scan(
			&address.ID,
			&address.FirstName,
			&address.LastName,
			&address.Email,
			&address.Mobile,
			&address.Address,
			&address.Country,
			&address.City,
			&address.State,
			&address.ZipCode,
			&address.IsDefault)

		if err != nil {
			return nil, err
		}
		Addresses = append(Addresses, address)
	}

	return Addresses, nil
}

func (m *DBModel) GetAllFeatureProducts(limitProduct int) ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// thực hiện query cho feature products
	rows, err := m.DB.QueryContext(ctx, `
	select id, name, image_path, stars, price, old_price, created_at
	from Product
	order by stars DESC
	limit ?
	`, limitProduct)

	// trả về error
	if err != nil {
		return nil, err
	}

	var featureProducts []Product
	// convert các products trong rows thành type Product
	// và lưu trong array
	for rows.Next() {
		var product Product
		err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.ImagePath,
			&product.Stars,
			&product.Price,
			&product.OldPrice,
			&product.CreatedAt)

		if err != nil {
			return nil, err
		}

		featureProducts = append(featureProducts, product)
	}

	return featureProducts, nil
}

func (m *DBModel) GetAllCategories() ([]Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// thực hiện query
	rows, err := m.DB.QueryContext(ctx, `
	select id, name, image_path
	from Category
	`)

	// trả về nill nếu xảy ra lỗi
	if err != nil {
		return nil, err
	}

	var categoryArray []Category
	for rows.Next() {
		var category Category
		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.ImagePath)

		if err != nil {
			return nil, err
		}
		categoryArray = append(categoryArray, category)
	}

	return categoryArray, nil
}

func (m *DBModel) GetRecentProducts() ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// do query in database
	rows, err := m.DB.QueryContext(ctx, `
	select id, name, image_path, stars, price, old_price, created_at
	from Product
	order by created_at DESC
	limit ?
	`, 10)

	// if error occurs, return nil
	if err != nil {
		return nil, err
	}

	// convert infor in rows to Product type
	var RecentProducts []Product
	for rows.Next() {
		var product Product
		err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.ImagePath,
			&product.Stars,
			&product.Price,
			&product.OldPrice,
			&product.CreatedAt)
		// if error occurs, return nil
		if err != nil {
			return nil, err
		}
		RecentProducts = append(RecentProducts, product)
	}

	return RecentProducts, nil
}

func (m *DBModel) GetAllBrands() ([]Brand, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// do query
	rows, err := m.DB.QueryContext(ctx, `
	select id, name, image_path
	from brand
	`)

	// check if error
	if err != nil {
		return nil, err
	}

	// convert rows into list of brands
	var Brands []Brand
	for rows.Next() {
		var brand Brand
		err = rows.Scan(
			&brand.ID,
			&brand.Name,
			&brand.ImagePath,
		)

		if err != nil {
			return nil, err
		}

		Brands = append(Brands, brand)
	}

	return Brands, nil
}

func (m *DBModel) GetInforCategoriesWithNumofProducts() ([]Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// do query
	rows, err := m.DB.QueryContext(ctx, `
	select t1.id, t1.name, t1.image_path, count(t2.id) as num_of_product
	from category as t1
	left join product as t2 on t1.id = t2.category_id
	group by t1.id
	`)

	// check for error
	if err != nil {
		return nil, err
	}

	// convert rows into Category
	var Categories []Category
	for rows.Next() {
		var category Category
		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.ImagePath,
			&category.NumOfProduct,
		)
		if err != nil {
			return nil, err
		}

		Categories = append(Categories, category)
	}
	// return list of Categories
	return Categories, nil
}

func (m *DBModel) GetInfoBrandsWithNumofProduct() ([]Brand, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// do query
	rows, err := m.DB.QueryContext(ctx, `
	select t1.id, t1.name, t1.image_path, count(t2.id) as num_of_product
	from brand as t1
	left join product as t2
	on t1.id = t2.brand_id
	group by 
	t1.id
	`)

	// check if error occurs
	if err != nil {
		return nil, err
	}

	// convert rows into Brand
	var Brands []Brand
	for rows.Next() {
		var brand Brand
		err = rows.Scan(
			&brand.ID,
			&brand.Name,
			&brand.ImagePath,
			&brand.NumOfProduct,
		)

		if err != nil {
			return nil, err
		}

		Brands = append(Brands, brand)
	}

	return Brands, nil
}

func (m *DBModel) GetInforImagesWithProductId(product_id int) ([]Image, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, `
	select id, name, image_path
	from image as  t1
	where t1.product_id = ?
	`, product_id)

	if err != nil {
		return nil, err
	}

	var Images []Image
	for rows.Next() {
		var image Image
		err = rows.Scan(&image.ID, &image.Name, &image.ImagePath)
		if err != nil {
			return nil, err
		}

		Images = append(Images, image)
	}

	return Images, nil
}

func (m *DBModel) GetBoughtProductWithUserId(userId int) ([]Invoice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// do query
	rows, err := m.DB.QueryContext(ctx, `
	select t2.name, t2.image_path, t1.quantity, t1.price, t1.status, t1.updated_at as bought_date
	from invoice as t1
	inner join product as t2 
	where t1.product_id = t2.id and t1.user_id = ?
	`, userId)

	if err != nil {
		return nil, err
	}

	// convert result of sql into ListInvoice object
	var ListInvoice []Invoice
	for rows.Next() {
		var invoice Invoice
		err = rows.Scan(&invoice.ProductName, &invoice.ImagePath, &invoice.Quantity, &invoice.Price, &invoice.Status, &invoice.BoughtDate)
		if err != nil {
			return nil, err
		}
		ListInvoice = append(ListInvoice, invoice)
	}

	return ListInvoice, nil
}

func (m *DBModel) GetProductsWithFilter(sort string, search string, tagId int, categoryId int, brandId int, offset int, limit int) ([]Product, error) {
	queryInfor := []interface{}{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "select t1.id, t1.name, t1.image_path, t1.stars, t1.price, t1.old_price from product t1 "

	// add query tagId
	if tagId != 0 {
		stmt += "inner join producttag t2 on t1.id = t2.product_id where t2.tag_id = ? "
		queryInfor = append(queryInfor, tagId)
	}

	// add query Search Keyword
	if search != "" {
		if strings.Contains(stmt, "where") {
			stmt += "and t1.name like ? "
		} else {
			stmt += "where t1.name like ? "
		}
		queryInfor = append(queryInfor, "%"+search+"%")
	}

	// add query categoryId
	if categoryId != 0 {
		if strings.Contains(stmt, "where") {
			stmt += "and t1.category_id = ? "
		} else {
			stmt += "where t1.category_id = ? "
		}
		queryInfor = append(queryInfor, categoryId)
	}

	// add query brandId
	if brandId != 0 {
		if strings.Contains(stmt, "where") {
			stmt += "and t1.brand_id = ? "
		} else {
			stmt += "where t1.brand_id = ? "
		}
		queryInfor = append(queryInfor, brandId)
	}

	// add query sort
	stmt += "order by "
	if sort == "price" {
		stmt += "t1.price DESC "
	} else if sort == "newest" {
		stmt += "t1.created_at DESC "
	} else if sort == "popular" {
		stmt += "t1.stars DESC "
	}

	// add query page
	if limit != 0 {
		stmt += "limit ? offset ? "
		queryInfor = append(queryInfor, limit)
		queryInfor = append(queryInfor, offset)
	}

	var ProductList []Product
	rows, err := m.DB.QueryContext(ctx, stmt, queryInfor...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var product Product
		err = rows.Scan(&product.ID, &product.Name, &product.ImagePath, &product.Stars, &product.Price, &product.OldPrice)
		if err != nil {
			return nil, err
		}
		ProductList = append(ProductList, product)
	}

	return ProductList, nil
}

func (m *DBModel) GetRelatedProductsInfor(productId int) ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, `
	select t1.id, t1.name,  t1.image_path, t1.old_price, t1.price
	from product as t1
	inner join producttag as t2
	where t1.id = t2.product_id and t1.id != ? and t2.tag_id in (
		select tag_id
		from producttag
		where product_id = ?
	)`, productId, productId)

	if err != nil {
		return nil, err
	}

	var RelatedProducts []Product
	for rows.Next() {
		var product Product
		err = rows.Scan(&product.ID, &product.Name, &product.ImagePath, &product.OldPrice, &product.Price)
		if err != nil {
			return nil, err
		}
		RelatedProducts = append(RelatedProducts, product)
	}

	return RelatedProducts, nil
}

/*
 * hàm tính toán số loại sản phẩm trong giỏ hàng của người dùng
 * nếu không có sản phẩm thì trả về giá trị 0
 */
func (m *DBModel) CountNumofTypeProductInCartByUserEmail(userEmail string) (int, error) {
	// nếu người dùng chưa đăng nhập
	if userEmail == "" {
		return 0, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// lấy thông tin user
	user, err := m.GetUserByEmail(userEmail)
	if err != nil {
		return 0, err
	}

	// chứa số loại sản phẩm trong giỏ hàng
	var quantity int

	// xử lý query
	err = m.DB.QueryRowContext(ctx, `
	select COUNT(*) as quantity
	from CartDetail
	where user_id = ?
	`, user.ID).Scan(&quantity)

	// xảy ra lỗi, quantity mặc định là 0
	if err != nil {
		return 0, err
	}

	// trả kết quả
	return quantity, nil
}
