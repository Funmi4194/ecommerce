package order

import (
	"github.com/funmi4194/ecommerce/enum"
	"github.com/funmi4194/ecommerce/repository/common"
	"github.com/uptrace/bun"
)

type Order struct {
	bun.BaseModel `bun:"table:orders" rsf:"false"`

	ID     string `bun:"id,pk" json:"id"`
	UserID string `bun:"user_id" json:"user_id"`

	Status      enum.OrderStatus `bun:"status" json:"status"`
	Reference   string           `bun:"reference" json:"reference"`
	Paid        bool             `bun:"paid" json:"paid"`
	PaidAt      bun.NullTime     `bun:"paid_at" json:"paid_at"`
	Cancelled   bool             `bun:"cancelled" json:"cancelled"`
	CancelledAt bun.NullTime     `bun:"cancelled_at" json:"cancelled_at"`
	Failed      bool             `bun:"failed" json:"failed"`
	FailedAt    bun.NullTime     `bun:"failed_at" json:"failed_at"`

	// SHA256 of the transaction.Invoice (prevents tampering & duplication)
	Checksum string `bun:"checksum" json:"checksum"`

	History []common.History `bun:"history,type:jsonb" json:"history" rsfr:"false"`

	// items that make up the total amount
	Invoice []Item `bun:"invoice,type:jsonb" json:"invoice" rsfr:"false"`

	// the total amount to be paid (including all possible fees)
	Amount float64 `bun:"amount" json:"amount"`

	Remark string `bun:"remark" json:"remark"`

	ProductID string `bun:"product_id" json:"product_id"`

	CreatedAt bun.NullTime `bun:"created_at" json:"created_at" rsfr:"false"`
	UpdatedAt bun.NullTime `bun:"updated_at" json:"updated_at" rsfr:"false"`
}

// schematic representation of an item in a order's invoice
type Item struct {
	// the key of the item
	Key string `json:"key"`

	// the name of the item
	Name string `json:"name"`

	// the amount of the item
	Amount float64 `json:"amount"`

	// the quantity of the item
	Quantity int `json:"quantity"`

	// the metadata for the item (can also be a json string)
	Metadata string `json:"metadata"`
}

type Orders []Order
