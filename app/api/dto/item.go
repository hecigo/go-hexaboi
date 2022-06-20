package dto

import (
	"hoangphuc.tech/go-hexaboi/domain/model"
)

type Item struct {
	model.Item
}

// Used to validate on creation
type ItemCreated struct {
	Code              string            `json:"code" validate:"required,max=50"`
	Name              string            `json:"name" validate:"required"`
	ShortName         string            `json:"short_name"`
	VariantAttributes map[string]string `json:"variant_attributes"`
	MasterSKU         string            `json:"master_sku"`
	PrimaryCategoryID uint              `json:"primary_category_id" validate:"gte=1"`
	CategoriesID      []uint            `json:"categories_id"`
	BrandID           uint              `json:"brand_id" validate:"gte=1"`
	CreatedBy         string            `json:"created_by" validate:"required"`
	UpdatedBy         string            `json:"updated_by" validate:"required"`
}

func (ic *ItemCreated) ToModel() *model.Item {
	mi := &model.Item{
		Code: ic.Code,
		Name: ic.Name,
	}

	// for _, c := range ic.CategoriesID {
	// 	mi.Categories = append(mi.Categories, &model.Category{
	// 		EntityID: base.EntityID{
	// 			ID: c,
	// 		},
	// 	})
	// }

	return mi
}

type ItemUpdated struct {
	Name              string            `json:"name"`
	ShortName         string            `json:"short_name"`
	VariantAttributes map[string]string `json:"variant_attributes"`
	MasterSKU         string            `json:"master_sku"`
	PrimaryCategoryID uint              `json:"primary_category_id" validate:"omitempty,gt=0"`
	CategoriesID      []uint            `json:"categories_id" validate:"dive,gt=0"`
	BrandID           uint              `json:"brand_id" validate:"omitempty,gt=0"`
	UpdatedBy         string            `json:"updated_by" validate:"required"`
}

type ItemFilterDto struct {
	Items     []string `query:"items" validate:"required"`
	PageIndex int      `query:"page_index" validate:"required"`
	PageSize  int      `query:"page_size" validate:"min=0,max=32"`
}
