package model

import "fmt"

// Item grouping, category
type Inventory struct {
	Code      string  `json:"inventory_code"`
	StoreCode string  `json:"store_code"`
	ItemCode  string  `json:"item_code"`
	Stock     float32 `json:"stock"`

	IsActived int8 `json:"is_actived"`
	IsDeleted int8 `json:"is_deleted"`
}

func (i Inventory) String() string {
	return i.StoreCode + "\t" + i.ItemCode + "\t" + fmt.Sprintf("%f", i.Stock)
}
