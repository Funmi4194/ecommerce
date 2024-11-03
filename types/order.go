package types

import (
	"time"

	"github.com/funmi4194/ecommerce/enum"
)

type InitiateOrder struct {
	Items []Item `json:"items"`
}

type Item struct {
	ProductId string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type ListOrder struct {
	UserId string           `json:"userId"`
	Status enum.OrderStatus `json:"status"`
}

type UpdateOrder struct {
	OrderId string           `json:"order_id"`
	Status  enum.OrderStatus `json:"status"`
}

type CancelOrder struct {
	OrderId string `json:"order_id"`
	Cancel  bool   `json:"cancel"`
}

type OrderFilter struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	MinAmount *float64  `json:"min_amount"`
	MaxAmount *float64  `json:"max_amount"`

	OrderId   string `json:"order_id"`
	Reference string `json:"reference"`
	Paid      *bool  `json:"paid"`
	Failed    *bool  `json:"failed"`
	Cancelled *bool  `json:"cancelled"`
	UserId    string `json:"user_id"`

	// pagination
	Page  *int `json:"page"`
	Limit *int `json:"limit"`
	// when true, the response will contain the pagination metadata
	Paginate bool `json:"paginate"`
}
