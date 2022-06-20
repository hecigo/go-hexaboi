package orm

import "hoangphuc.tech/go-hexaboi/domain/model"

// Items can be goods, products, gifts, services...
type Item struct {
	Code         string `json:"item_code"`
	Name         string `json:"item_name"`
	Description  string `json:"item_description"`
	BrandCode    string `json:"brand_code"`
	Image        string `json:"item_image"`
	CategoryCode string `json:"category_code"`
	ListedPrice  int32  `json:"listed_price"`
	SalesPrice   int32  `json:"sales_price"`
	FullPrice    int32  `json:"full_price"`
	Type         string `json:"type"`
	IsActived    int8   `json:"is_actived"`
	IsDeleted    int8   `json:"is_deleted"`

	// The key-value collection of variant attributes. Ex: color, size...
	// ItemAttribute *map[string]string `json:"item_attribute"`
	ItemAttribute string `json:"item_attribute"`

	// This SKU represents a group of SKU of the same type.
	// It is the parent/configurable SKU as well.
	ParentCode string `json:"parent_code"`

	// All categories including SKU
	Categories []*Category `json:"categorys"`

	// Brand
	Brands []*Brand `json:"brands"`

	// danh sách tồn kho của 1 sản phẩm
	Inventorys []*Inventory `json:"inventorys"`
}

func (o *Item) ToModel(m *model.Item) {
	if m == nil {
		m = new(model.Item)
	}
	m.Code = o.Code
	m.Name = o.Name
	m.ParentCode = o.ParentCode

	// o.Inventorys.ToModel(&m.Inventorys)

	// o.Brands.ToModel(&m.Brands)
	// if m.Brands.Code == "" {
	// 	m.Brands.Code = o.BrandCode
	// }

	m.ItemAttribute = o.ItemAttribute

	if o.Inventorys != nil && len(o.Inventorys) > 0 {
		m.Inventorys = make([]*model.Inventory, len(o.Inventorys))
		for i, oc := range o.Inventorys {
			m.Inventorys[i] = new(model.Inventory)
			oc.ToModel(m.Inventorys[i])
		}
	}

	if o.Brands != nil && len(o.Brands) > 0 {
		m.Brands = make([]*model.Brand, len(o.Brands))
		for i, oc := range o.Brands {
			m.Brands[i] = new(model.Brand)
			oc.ToModel(m.Brands[i])
		}
	}

	if o.Categories != nil && len(o.Categories) > 0 {
		m.Categories = make([]*model.Category, len(o.Categories))
		for i, oc := range o.Categories {
			m.Categories[i] = new(model.Category)
			oc.ToModel(m.Categories[i])
		}
	}
}
