package models

import (
	"database/sql"
	"time"
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

type OrderDetail struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	ImagePath string  `json:"image_path"`
	Quantity  int     `json:"quantity"`
	Price     float32 `json:"price"`
	Total     float32 `json:"total"`
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
