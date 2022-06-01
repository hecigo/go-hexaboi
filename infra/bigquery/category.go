package bigquery

import (
	"fmt"

	"hoangphuc.tech/hercules/domain/base"
	"hoangphuc.tech/hercules/domain/model"
	"hoangphuc.tech/hercules/infra/orm"
)

type CategoryRepository struct {
	base.Repository
}

func (r *CategoryRepository) FindAll() ([]orm.Category, error) {
	var bqResult []map[string]interface{}
	result := DB().Table("category").Select("code", "description").Find(&bqResult)
	var categories []orm.Category
	for _, bqCate := range bqResult {
		b := orm.Category{
			Code:       fmt.Sprintf("%v", bqCate["code"]),
			Name:       fmt.Sprintf("%v", bqCate["description"]),
			DivisionBy: model.DIVISION_MERCHANDISE,
			Entity: orm.Entity{
				CreatedBy: "system",
				UpdatedBy: "system",
			},
		}
		categories = append(categories, b)
	}
	return categories, result.Error
}
