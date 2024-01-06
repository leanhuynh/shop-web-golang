package models

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
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

func (m *DBModel) CreateUser(email string, password string, firstname string, lastname string, mobile string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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

func (m *DBModel) InsertOrder(token string, quantity int, price float32) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var user *User
	user, err := m.GetUserForToken(token)
	if err != nil {
		return false, errors.New("no match token")
	}

	_, err = m.DB.ExecContext(ctx, `insert into PlaceOrder (user_id, quantity, price) values (?, ?, ?)`, user.ID, quantity, price)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *DBModel) GetProductList(stmt string, data []interface{}) ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// execute query
	rows, err := m.DB.QueryContext(ctx, stmt, data...)
	if err != nil {
		return nil, err
	}

	// scan rows
	var ProductList []Product
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

func (m *DBModel) GetAllCategories() ([]Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		select 
			t1.id, t1.name, t1.image_path, count(t2.id) as num_of_product  
		from 
			Category as t1
		left join 
			product as t2 on t1.id = t2.category_id 
		group by 
			t1.id;
	`
	rows, err := m.DB.QueryContext(ctx, stmt)

	if err != nil {
		return nil, err
	}

	var Categories []Category

	for rows.Next() {
		var category Category
		err = rows.Scan(&category.ID, &category.Name, &category.ImagePath, &category.NumOfProduct)
		if err != nil {
			return nil, err
		}
		Categories = append(Categories, category)
	}

	return Categories, nil
}

func (m *DBModel) GetAllBrand() ([]Brand, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, `
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
		return nil, err
	}

	var Brands []Brand

	for rows.Next() {
		var brand Brand
		err = rows.Scan(&brand.ID, &brand.Name, &brand.ImagePath, &brand.NumOfProduct)
		if err != nil {
			return nil, err
		}
		Brands = append(Brands, brand)
	}

	return Brands, nil
}

func (m *DBModel) GetAllTag() ([]Tag, error) {
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
