package main

import "myapp/internal/models"

type Cart struct {
	Items           []models.OrderDetail `json:"items"`
	Quantity        int                  `json:"quantity"`
	SubTotal        float32              `json:"sub_total"`
	Total           float32              `json:"total"`
	Shipping        float32              `json:"shipping"`
	Discount        float32              `json:"discount"`
	CouponCode      string               `json:"coupon_code"`
	PaymentMethod   string               `json:"payment_method"`
	ShippingAddress string               `json:"shipping_address"`
}

func (cart *Cart) GetSubTotal() float32 {
	var subTotal float32
	for _, item := range cart.Items {
		subTotal += item.Total
	}

	return subTotal
}

func (cart *Cart) GetQuantity() int {
	var quantity int
	for _, item := range cart.Items {
		quantity += item.Quantity
	}

	return quantity
}

func (cart *Cart) GetTotal() float32 {
	return cart.GetSubTotal() + cart.Shipping - cart.Discount
}

func (cart *Cart) GetCart() Cart {
	return Cart{
		Items:           cart.Items,
		Quantity:        cart.GetQuantity(),
		SubTotal:        cart.GetSubTotal(),
		Total:           cart.GetTotal(),
		Shipping:        cart.Shipping,
		Discount:        cart.Discount,
		CouponCode:      cart.CouponCode,
		PaymentMethod:   cart.PaymentMethod,
		ShippingAddress: cart.ShippingAddress,
	}
}

// func NewCart(oldCart *Cart) *Cart {
// 	if oldCart == nil {
// 		return &Cart{
// 			Items:           nil,
// 			Shipping:        0,
// 			Discount:        0,
// 			CouponCode:      "",
// 			PaymentMethod:   "",
// 			ShippingAddress: "",
// 		}
// 	}
// }

// func (cart *Cart) GetQuantity() int {
// 	quantity := 0
// 	for id := range cart.Items {
// 		quantity += int(cart.Items[id].quantity)
// 	}

// 	return quantity
// }
