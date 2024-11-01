package product

import (
	"github.com/funmi4194/ecommerce/enum"
	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:products" rsf:"false"`
	ID            string             `bun:"id,pk" json:"id"`
	Name          string             `bun:"name" json:"name"`
	Price         float64            `bun:"price" json:"price"`
	Stock         int64              `bun:"stock" json:"stock"`
	ProductUrl    string             `bun:"product_url" json:"product_url"`
	Status        enum.ProductStatus `bun:"status" json:"status" rsfr:"false"`
	Description   string             `bun:"description" json:"description"`
	CreatedAt     bun.NullTime       `bun:"created_at" json:"created_at" rsfr:"false"`
	UpdatedAt     bun.NullTime       `bun:"updated_at" json:"updated_at" rsfr:"false"`
}

type Products []Product
