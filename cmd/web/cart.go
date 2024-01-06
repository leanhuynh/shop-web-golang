package main

// import (
// 	"myapp/internal/models"
// )

// type Cart struct {
// 	Items           []models.OrderDetail `json:"items"`
// 	Quantity        int                  `json:"quantity"`
// 	SubTotal        float32              `json:"sub_total"`
// 	Total           float32              `json:"total"`
// 	Shipping        float32              `json:"shipping"`
// 	Discount        float32              `json:"discount"`
// 	CouponCode      string               `json:"coupon_code"`
// 	PaymentMethod   string               `json:"payment_method"`
// 	ShippingAddress string               `json:"shipping_address"`
// }

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

// func (cart *Cart) AddCart(product models.Product, quantity int) {
// 	// kiem tra san pham co trong gio hang
// 	// neu co
// 	index := cart.GetIndexOfProductId(product.ID)
// 	if index != -1 {
// 		cart.Items[index].Quantity += quantity
// 		cart.Items[index].Price = product.Price
// 		cart.Items[index].Total = float32(cart.Items[index].Quantity) * product.Price
// 	} else { // neu san pham chua co trong gio hang
// 		cart.Items = append(cart.Items, models.OrderDetail{
// 			ID:       product.ID,
// 			Quantity: quantity,
// 			Price:    product.Price,
// 			Total:    product.Price * float32(quantity),
// 		})
// 	}
// }

// func (cart *Cart) UpdateCart(product models.Product, quantity int) {
// 	// kiem tra san pham co trong gio hang
// 	// neu co
// 	index := cart.GetIndexOfProductId(product.ID)
// 	if index != -1 {
// 		cart.Items[index].Quantity = quantity
// 		cart.Items[index].Price = product.Price
// 		cart.Items[index].Total = float32(cart.Items[index].Quantity) * product.Price
// 	} else { // neu san pham chua co trong gio hang
// 		cart.Items = append(cart.Items, models.OrderDetail{
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
// 		newArray := make([]models.OrderDetail, 0, len(cart.Items)-1)
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

// // func NewCart(oldCart *Cart) *Cart {
// // 	if oldCart == nil {
// // 		return &Cart{
// // 			Items:           nil,
// // 			Shipping:        0,
// // 			Discount:        0,
// // 			CouponCode:      "",
// // 			PaymentMethod:   "",
// // 			ShippingAddress: "",
// // 		}
// // 	}
// // }

// // func (cart *Cart) GetQuantity() int {
// // 	quantity := 0
// // 	for id := range cart.Items {
// // 		quantity += int(cart.Items[id].quantity)
// // 	}

// // 	return quantity
// // }
