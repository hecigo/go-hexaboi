package orm

import (
	model "hoangphuc.tech/hercules/domain/model"
	"hoangphuc.tech/hercules/infra/orm/extensions"
)

type Item struct {
	EntityID
	Code              string          `json:"code" gorm:"uniqueIndex; not null; type:varchar(50); check:code <> ''"`
	Name              string          `json:"name" gorm:"not null; check:name <> ''"`
	VariantAttributes extensions.JSON `json:"variant_attributes" gorm:"type:json"`
	MasterSKU         string          `json:"master_sku" gorm:"type:varchar(50)"`
	PrimaryCategoryID uint            `json:"primary_category_id"`
	PrimaryCategory   Category        `json:"primary_category" gorm:"foreignKey:PrimaryCategoryID; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	BrandID           uint            `json:"brand_id"`
	Brand             Brand           `json:"brand" gorm:"foreignKey:BrandID; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Categories        []Category      `json:"categories" gorm:"many2many:item_j_category; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Entity
}

// Item constructor
func NewItem(_item *model.Item) *Item {
	item := &Item{
		EntityID:  *NewEntityID(_item.ID),
		Code:      _item.Code,
		Name:      _item.Name,
		MasterSKU: _item.MasterSKU,
		Entity:    *NewEntity(_item.Entity),
	}

	item.VariantAttributes = make(extensions.JSON)
	item.VariantAttributes.Marshal(_item.VariantAttributes)

	// Brand
	item.BrandID = _item.Brand.ID
	item.Brand = *NewBrand(_item.Brand)

	// Categories
	item.PrimaryCategoryID = _item.PrimaryCategory.ID
	item.PrimaryCategory = *NewCategory(_item.PrimaryCategory)

	for _, c := range _item.Categories {
		item.Categories = append(item.Categories, *NewCategory(c))
	}

	return item
}

func (o *Item) ToModel(m *model.Item) {
	if m == nil {
		m = new(model.Item)
	}
	o.Entity.ToModel(&m.Entity)
	m.ID = o.ID
	m.Code = o.Code
	m.Name = o.Name
	m.MasterSKU = o.MasterSKU
	o.Brand.ToModel(&m.Brand)
	o.PrimaryCategory.ToModel(&m.PrimaryCategory)
	m.VariantAttributes = o.VariantAttributes.Unmarshal()
	if o.Categories != nil && len(o.Categories) > 0 {
		m.Categories = make([]model.Category, len(o.Categories))
		for i, c := range o.Categories {
			c.ToModel(&m.Categories[i])
		}
	}
}
