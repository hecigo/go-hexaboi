package bigquery

import (
	"fmt"

	"hecigo.com/go-hexaboi/domain/base"
	"hecigo.com/go-hexaboi/infra/orm"
)

type BrandRepository struct {
	base.Repository
}

func (r *BrandRepository) FindAll() ([]orm.Brand, error) {
	var bqResult []map[string]interface{}
	result := DB().Table("brand").Select("code", "description").Find(&bqResult)
	var brands []orm.Brand
	for _, bqBrand := range bqResult {
		b := orm.Brand{
			Code: fmt.Sprintf("%v", bqBrand["code"]),
			Name: fmt.Sprintf("%v", bqBrand["description"]),
			Entity: orm.Entity{
				CreatedBy: "system",
				UpdatedBy: "system",
			},
		}
		brands = append(brands, b)
	}
	return brands, result.Error
}
