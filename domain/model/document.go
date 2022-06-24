package model

import "hoangphuc.tech/go-hexaboi/domain/base"

type Document struct {
	base.EntityID
	Code          string       `json:"code"`
	Type          DocumentType `json:"type"`
	SumOfQuantity uint         `json:"sum_qty"`
	FromWarehouse Warehouse    `json:"from_wh"`
	ToWarehouse   Warehouse    `json:"to_wh"`
	base.Entity
}

type DocumentType string

const (
	SALES_ORDER    = "SO"
	PURCHASE_ORDER = "PO"
	TRANSFER_ORDER = "TO"
)

type DocumentItem struct {
	base.EntityID
	Document Document `json:"document"`
	ItemCode string   `json:"item_code"`
	Quantity uint     `json:"qty"`
}
