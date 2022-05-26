package orm

import (
	"errors"

	"gorm.io/gorm"
	model "hoangphuc.tech/hercules/domain/model"
	"hoangphuc.tech/hercules/infra/orm/ext"
)

type Item struct {
	EntityID
	Code              string     `json:"code" gorm:"uniqueIndex; not null; type:varchar(50); check:code <> ''"`
	Name              string     `json:"name" gorm:"not null; check:name <> ''"`
	ShortName         string     `json:"short_name"`
	VariantAttributes ext.JSON   `json:"variant_attributes"`
	MasterSKU         string     `json:"master_sku" gorm:"type:varchar(50)"`
	PrimaryCategoryID uint       `json:"primary_category_id" gorm:"not null;"`
	PrimaryCategory   Category   `json:"primary_category" gorm:"foreignKey:PrimaryCategoryID; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	BrandID           uint       `json:"brand_id" gorm:"not null;"`
	Brand             Brand      `json:"brand" gorm:"foreignKey:BrandID; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Categories        []Category `json:"categories" gorm:"many2many:item_j_category; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Entity
}

// Initialize orm.Category from model.Category
func NewItem(item *model.Item) *Item {
	result := &Item{
		EntityID:  *NewEntityID(item.ID),
		Code:      item.Code,
		Name:      item.Name,
		ShortName: item.ShortName,
		MasterSKU: item.MasterSKU,
		Entity:    *NewEntity(item.Entity),
	}

	if item.VariantAttributes != nil {
		result.VariantAttributes = *new(ext.JSON)
		result.VariantAttributes.Load(item.VariantAttributes)
	}

	// Brand
	result.BrandID = item.Brand.ID
	result.Brand = *NewBrand(&item.Brand)

	// Categories
	result.PrimaryCategoryID = item.PrimaryCategory.ID
	result.PrimaryCategory = *NewCategory(&item.PrimaryCategory)

	for _, c := range item.Categories {
		result.Categories = append(result.Categories, *NewCategory(&c))
	}

	return result
}

// Scan orm.Item into model.Item
func (o *Item) ToModel(m *model.Item) {
	if m == nil {
		m = new(model.Item)
	}
	o.Entity.ToModel(&m.Entity)
	m.ID = o.ID
	m.Code = o.Code
	m.Name = o.Name
	m.ShortName = o.ShortName
	m.MasterSKU = o.MasterSKU
	o.Brand.ToModel(&m.Brand)
	o.PrimaryCategory.ToModel(&m.PrimaryCategory)

	variants, err := o.VariantAttributes.ToStrMap()
	if err != nil {
		panic(err)
	}
	m.VariantAttributes = variants

	if o.Categories != nil && len(o.Categories) > 0 {
		m.Categories = make([]model.Category, len(o.Categories))
		for i, c := range o.Categories {
			c.ToModel(&m.Categories[i])
		}
	}
}

func (u *Item) BeforeUpdate(tx *gorm.DB) (err error) {
	// if Role changed
	if tx.Statement.Changed("VariantAttributes") {
		return errors.New("VariantAttributes not allowed to change")
	}

	return nil
}
