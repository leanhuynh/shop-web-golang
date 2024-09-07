package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Cart struct {
	UserId     int          `json:"user_id"`
	CartDetail []CartDetail `json:"cart_detail"`
	Quantity   int          `json:"quantity"`
	SubTotal   float64      `json:"sub_total"`
	Total      float64      `json:"total"`
	// Shipping        float32       `json:"shipping"`
	// Discount        float32       `json:"discount"`
	// CouponCode      string        `json:"coupon_code"`
	// PaymentMethod   string        `json:"payment_method"`
	// ShippingAddress string        `json:"shipping_address"`
}

type CartDetail struct {
	UserId    int     `json:"user_id"`
	ProductId int     `json:"product_id"`
	Name      string  `json:"name"`
	ImagePath string  `json:"image_path"`
	Quantity  int     `json:"quantity"`
	Price     float32 `json:"price"`
	Total     float32 `json:"total"`
}

// func (cart *Cart) GetSubTotal() float32 {
// 	var subTotal float32
// 	for _, item := range cart.Items {
// 		subTotal += item.Total
// 	}

// 	return subTotal
// }

// func (cart *Cart) GetQuantity() int {
// 	var quantity int
// 	for _, item := range cart.Items {
// 		quantity += item.Quantity
// 	}

// 	return quantity
// }

// func (cart *Cart) GetTotal() float32 {
// 	return cart.GetSubTotal() + cart.Shipping - cart.Discount
// }

// func (cart *Cart) GetCart() Cart {
// 	return Cart{
// 		Items:           cart.Items,
// 		Quantity:        cart.GetQuantity(),
// 		SubTotal:        cart.GetSubTotal(),
// 		Total:           cart.GetTotal(),
// 		Shipping:        cart.Shipping,
// 		Discount:        cart.Discount,
// 		CouponCode:      cart.CouponCode,
// 		PaymentMethod:   cart.PaymentMethod,
// 		ShippingAddress: cart.ShippingAddress,
// 	}
// }

// func (cart *Cart) AddCart(product Product, quantity int) {
// 	// kiem tra san pham co trong gio hang
// 	// neu co
// 	index := cart.GetIndexOfProductId(product.ID)
// 	if index != -1 {
// 		cart.Items[index].Quantity += quantity
// 		cart.Items[index].Price = product.Price
// 		cart.Items[index].Total = float32(cart.Items[index].Quantity) * product.Price
// 	} else { // neu san pham chua co trong gio hang
// 		cart.Items = append(cart.Items, OrderDetail{
// 			ID:       product.ID,
// 			Quantity: quantity,
// 			Price:    product.Price,
// 			Total:    product.Price * float32(quantity),
// 		})
// 	}
// }

// func (cart *Cart) UpdateCart(product Product, quantity int) {
// 	// kiem tra san pham co trong gio hang
// 	// neu co
// 	index := cart.GetIndexOfProductId(product.ID)
// 	if index != -1 {
// 		cart.Items[index].Quantity = quantity
// 		cart.Items[index].Price = product.Price
// 		cart.Items[index].Total = float32(cart.Items[index].Quantity) * product.Price
// 	} else { // neu san pham chua co trong gio hang
// 		cart.Items = append(cart.Items, OrderDetail{
// 			ID:       product.ID,
// 			Quantity: quantity,
// 			Price:    product.Price,
// 			Total:    product.Price * float32(quantity),
// 		})
// 	}
// }

// func (cart *Cart) RemoveCart(product_id int) {
// 	// kiem tra san pham co trong gio hang
// 	// neu co
// 	index := cart.GetIndexOfProductId(product_id)
// 	// kiem tra gia tri index co hop le
// 	if index != -1 {
// 		newArray := make([]OrderDetail, 0, len(cart.Items)-1)
// 		newArray = append(newArray, cart.Items[0:index]...)
// 		newArray = append(newArray, cart.Items[index+1:]...)
// 		cart.Items = newArray
// 	}
// }

// func (cart *Cart) GetIndexOfProductId(product_id int) int {
// 	// kiem tra gia tri index co hop le
// 	for index, order := range cart.Items {
// 		// neu tim thay san pham phu hop voi given product_id
// 		if order.ID == product_id {
// 			return index
// 		}
// 	}

// 	// neu khong tim thay
// 	// tra ve -1
// 	return -1
// }

func (m *DBModel) GetCartForUser(user_id int) ([]CartDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `select t1.user_id, t1.product_id, t2.name, t2.image_path, t1.quantity, t1.price, t1.total 
	from 
		CartDetail t1
	inner 
		join Product t2
	on 
		t1.product_id = t2.id
	where 
		t1.user_id = ?`
	rows, err := m.DB.QueryContext(ctx, stmt, user_id)
	if err != nil {
		return nil, err
	}

	// scan rows
	var CartDetailList []CartDetail
	for rows.Next() {
		var cartDetail CartDetail
		err = rows.Scan(&cartDetail.UserId, &cartDetail.ProductId, &cartDetail.Name, &cartDetail.ImagePath, &cartDetail.Quantity, &cartDetail.Price, &cartDetail.Total)
		if err != nil {
			return nil, err
		}
		CartDetailList = append(CartDetailList, cartDetail)
	}

	return CartDetailList, nil
}

func (m *DBModel) UpdateProductToCart(user_id int, product Product, quantity int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// neu muon xoa voi quantity = 0
	if quantity == 0 {
		return m.RemoveProduct(user_id, product.ID)
	}

	// neu muon update voi so luong khac 0
	stmt := `update CartDetail set quantity = ? , price = ? , total = ? , updated_at = ?
		where 
			user_id = ? and product_id = ? 
	`
	_, err := m.DB.ExecContext(ctx, stmt, quantity, product.Price, float32(quantity)*product.Price, time.Now(), user_id, product.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) AddProducToCart(user_id int, product Product, quantity int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// truy vấn để tìm thông tin số lượng sản phẩm còn tổn tại trong cơ sở dữ liệu
	stmt := `select quantity from CartDetail where user_id = ? and product_id = ?`
	row := m.DB.QueryRowContext(ctx, stmt, user_id, product.ID)
	// var orderProduct Product

	// // lấy số lượng trong kho của sản phẩm
	// row := m.DB.QueryRowContext(ctx, `select quantity from product where id = ?`, product.ID)

	// // convert product
	// err := row.Scan(&orderProduct.Quantity)

	// // nếu có lỗi, trả về
	// if err != nil {
	// 	return err
	// }

	// // kiểm tra nếu số lượng yêu cầu có vượt quá số lượng trong kho hay không
	// if orderProduct.Quantity < quantity {
	// 	return errors.New("số lượng thêm vào vượt quá số lượng hiện có")
	// }

	if row != nil {
		// scan data
		var order struct {
			Quantity int `json:"quantity"`
		}
		err := row.Scan(&order.Quantity)

		// thực sự có lỗi trong scan
		if err != nil {
			// nếu lỗi sai không phải là do không có sản phẩm cũng loại trong giỏ hàng
			if err != sql.ErrNoRows {
				return err
			} else {
				if quantity > product.Quantity { // nếu số lượng muốn thêm vào > số lượng hiện có
					return errors.New("số lượng thêm vào vượt quá số lượng hiện có")
				}
			}
		} else {
			// nếu có sản phẩm cùng loại trong giỏ hàng
			fmt.Printf("so luong trong gio hang: %d\n", order.Quantity)
			// neu khong co loi
			if order.Quantity+quantity > product.Quantity { // nếu số lượng muốn thêm vào > số lượng hiện có
				return errors.New("số lượng thêm vào vượt quá số lượng hiện có")
			}
			return m.UpdateProductToCart(user_id, product, order.Quantity+quantity)
		}
	}

	// thực hiện truy vấn thêm sản phẩm vào giỏ hàng
	stmt = `insert into CartDetail (user_id, product_id, quantity, price, total)
	values (?, ?, ?, ?, ?)`
	_, err := m.DB.ExecContext(ctx, stmt, user_id, product.ID, quantity, product.Price, float32(quantity)*product.Price)

	// kiểm tra nếu có lỗi xảy ra, trả về error
	if err != nil {
		return err
	}

	// // cập nhật lại số lượng trong kho của sản phẩm
	// _, err = m.DB.ExecContext(ctx, `
	// 	update product set quantity = ?
	// 	where id = ?
	// `, orderProduct.Quantity-quantity, product.ID)

	// // nếu có lỗi xảy ra, trả về error
	// if err != nil {
	// 	return err
	// }

	// thêm sản phẩm thành công vào giỏ hàng, trả về nil
	return nil
}

func (m *DBModel) RemoveProduct(user_id int, product_id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := ` delete from CartDetail 
	where 
		user_id = ? and product_id = ?
	`
	_, err := m.DB.ExecContext(ctx, stmt, user_id, product_id)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) RemoveCart(user_id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `delete from CartDetail
	where 
		user_id = ? 
	`
	_, err := m.DB.ExecContext(ctx, stmt, user_id)
	if err != nil {
		return err
	}

	return nil
}
