package orm

import (
	"fmt"

	"hoangphuc.tech/dora/domain/model"
)

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

// Scan orm.Brand into model.Brand
func (i *Inventory) ToModel(inven *model.Inventory) {
	if inven == nil {
		inven = new(model.Inventory)
	}
	inven.Code = i.Code
	inven.StoreCode = i.StoreCode
	inven.ItemCode = i.ItemCode
	inven.Stock = i.Stock
	inven.IsActived = i.IsActived
	inven.IsDeleted = i.IsDeleted
}
