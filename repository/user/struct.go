package user

import (
	"github.com/funmi4194/ecommerce/enum"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users" rsf:"false"`
	ID            string       `bun:"id,pk" json:"id"`
	Email         string       `bun:"email,unique" json:"email"`
	Password      string       `bun:"password" json:"password"`
	Role          enum.Role    `bun:"role" json:"role" rsfr:"false"`
	CreatedAt     bun.NullTime `bun:"created_at" json:"created_at" rsfr:"false"`
	UpdatedAt     bun.NullTime `bun:"updated_at" json:"updated_at" rsfr:"false"`
}
