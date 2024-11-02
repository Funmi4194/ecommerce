package types

import (
	"github.com/funmi4194/ecommerce/enum"
	"github.com/uptrace/bun"
)

type Publish struct {
	Products []Product `json:"products"`
}

type Delete struct {
	Products []Product `json:"products"`
}

type Product struct {
	ProductId   string             `json:"product_id"`
	Name        string             `json:"name"`
	Price       float64            `json:"price"`
	Stock       int64              `json:"stock"`
	ProductUrl  string             `json:"product_url"`
	Status      enum.ProductStatus `json:"status"`
	Description string             `json:"description"`
}

type UpdateProduct struct {
	ProductId   string              `json:"product_id"`
	Name        *string             `json:"name"`
	Price       *float64            `json:"price"`
	Stock       *int64              `json:"stock"`
	ProductUrl  *string             `json:"product_url"`
	Status      *enum.ProductStatus `json:"status"`
	Description *string             `json:"description"`
}

type ProductFilter struct {
	ProductId string             `json:"product_id"`
	MinAmount *int64             `json:"min_amount"`
	MaxAmount *int64             `json:"max_amount"`
	Status    enum.ProductStatus `json:"status"`

	StartDate bun.NullTime `json:"start_date"`
	EndDate   bun.NullTime `json:"end_date"`

	// searches on name, description
	Search string `json:"search"`

	// pagination
	Page  *int `json:"page"`
	Limit *int `json:"limit"`

	// when true, the response will contain the pagination metadata
	Paginate bool `json:"paginate"`
}
