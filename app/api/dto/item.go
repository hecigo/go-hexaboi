package dto

import (
	"hoangphuc.tech/hercules/domain/base"
	"hoangphuc.tech/hercules/domain/model"
)

type Item struct {
	model.Item
}

// Used to validate on creation
type ItemCreated struct {
	Code              string            `json:"code" validate:"required"`
	Name              string            `json:"name" validate:"required"`
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
		Code:              ic.Code,
		Name:              ic.Name,
		VariantAttributes: ic.VariantAttributes,
		MasterSKU:         ic.MasterSKU,
		PrimaryCategory: model.Category{
			EntityID: base.EntityID{
				ID: ic.PrimaryCategoryID,
			},
		},
		Brand: model.Brand{
			EntityID: base.EntityID{
				ID: ic.BrandID,
			},
		},
		Entity: base.Entity{
			CreatedBy: ic.CreatedBy,
			UpdatedBy: ic.UpdatedBy,
		},
	}

	for _, c := range ic.CategoriesID {
		mi.Categories = append(mi.Categories, model.Category{
			EntityID: base.EntityID{
				ID: c,
			},
		})
	}

	return mi
}

type ItemFilterDto struct {
	Items     []string `query:"items" validate:"required"`
	PageIndex int      `query:"page_index" validate:"required"`
	PageSize  int      `query:"page_size" validate:"min=0,max=32"`
}
