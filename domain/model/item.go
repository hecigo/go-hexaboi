package model

// Items can be goods, products, gifts, services...
type Item struct {
	Code         string  `json:"item_code"`
	Name         string  `json:"item_name"`
	Description  string  `json:"item_description"`
	BrandCode    string  `json:"brand_code"`
	Image        string  `json:"item_image"`
	CategoryCode string  `json:"category_code"`
	ListedPrice  float32 `json:"listed_price"`
	SalesPrice   float32 `json:"sales_price"`
	FullPrice    float32 `json:"full_price"`
	Type         string  `json:"type"`
	IsActived    int8    `json:"is_actived"`
	IsDeleted    int8    `json:"is_deleted"`

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

func (b Item) String() string {
	return b.Code + "\t" + b.Name
}
