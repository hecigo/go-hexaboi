package orm

import (
	"database/sql"

	"gorm.io/gorm"
	model "hoangphuc.tech/hercules/domain/model"
	"hoangphuc.tech/hercules/infra/orm/ext"
)

type Item struct {
	EntityID
	Code                string      `json:"code" gorm:"uniqueIndex; not null; type:varchar(50);"`
	Name                string      `json:"name" gorm:"not null; check:name <> ''"`
	ShortName           string      `json:"short_name"`
	VariantAttributes   *ext.JSON   `json:"variant_attributes"`
	MasterSKU           string      `json:"master_sku" gorm:"type:varchar(50)"`
	PrimaryCategoryCode string      `json:"primary_category_code" gorm:"column:category_code; not null;"`
	PrimaryCategory     Category    `json:"primary_category" gorm:"foreignKey:PrimaryCategoryCode; references:Code; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	BrandCode           string      `json:"brand_code" gorm:"not null;"`
	Brand               Brand       `json:"brand" gorm:"foreignKey:BrandCode; references:Code; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Categories          []*Category `json:"categories" gorm:"many2many:item_j_category; foreignKey:Code; joinForeignKey:ItemCode; references: Code; joinReferences: CategoryCode; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
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
		result.VariantAttributes = new(ext.JSON)
		result.VariantAttributes.Scan(*item.VariantAttributes)
	}

	// Brand
	result.BrandCode = item.Brand.Code
	result.Brand = *NewBrand(&item.Brand)

	// Categories
	result.PrimaryCategoryCode = item.PrimaryCategory.Code
	result.PrimaryCategory = *NewCategory(&item.PrimaryCategory)

	if item.Categories != nil {
		for _, c := range item.Categories {
			result.Categories = append(result.Categories, NewCategory(c))
		}
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

	o.PrimaryCategory.ToModel(&m.PrimaryCategory)
	if m.PrimaryCategory.Code == "" {
		m.PrimaryCategory.Code = o.PrimaryCategoryCode
	}

	o.Brand.ToModel(&m.Brand)
	if m.Brand.Code == "" {
		m.Brand.Code = o.BrandCode
	}

	if o.VariantAttributes != nil {
		variants, err := o.VariantAttributes.ToStrMap()
		if err != nil {
			panic(err)
		}
		m.VariantAttributes = variants
	}

	if o.Categories != nil && len(o.Categories) > 0 {
		m.Categories = make([]*model.Category, len(o.Categories))
		for i, oc := range o.Categories {
			m.Categories[i] = new(model.Category)
			oc.ToModel(m.Categories[i])
		}
	}
}

// Updating data in same transaction
func (item *Item) BeforeUpdate(tx *gorm.DB) (err error) {
	if (item.Categories != nil) && len(item.Categories) > 0 {
		// item.Categories is a many2many association, so we need to delete it before updating.
		// On updating, GORM will re-create new associations automatically.
		tx.Statement.Exec("DELETE FROM item_j_category WHERE item_id = @item_id", sql.Named("item_id", item.ID))
	}
	return nil
}
