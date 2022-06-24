package model

import "hoangphuc.tech/dora/domain/base"

// Transaction presents stock moving between different warehouses
type Transaction struct {
	base.EntityID
	Code              string    `json:"code"` // Trigger code of transaction, ensuring a document with specific status just moves stock once time
	Origin            Warehouse `json:"origin"`
	Target            Warehouse `json:"target"`
	Quantity          uint      `json:"qty"`
	ReferenceDocument Document  `json:"ref_doc"`
	base.Entity
}
