package model

import "github.com/uptrace/bun"

type Product struct {
	bun.BaseModel `bun:"table:products,alias:p"`
	ID            uint64 `bun:",pk,autoincrement"`
	Name          string `bun:"name,notnull"`
	Description   string `bun:"description,notnull"`
	Price         uint16 `bun:"price,notnull"`
	TotalPrice    uint16 `bun:"totalPrice,notnull"`
	Count         uint   `bun:"count,notnull"`
}
