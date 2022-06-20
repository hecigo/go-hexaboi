package orm

import "hoangphuc.tech/go-hexaboi/domain/model"

type Brand struct {
	Code      string `json:"brand_code"`
	Name      string `json:"brand_name"`
	IsActived int8   `json:"is_actived"`
	IsDeleted int8   `json:"is_deleted"`
}

// Scan orm.Brand into model.Brand
func (b *Brand) ToModel(brand *model.Brand) {
	if brand == nil {
		brand = new(model.Brand)
	}
	brand.Code = b.Code
	brand.Name = b.Name
	brand.IsActived = b.IsActived
	brand.IsDeleted = b.IsDeleted
}
