package bigquery

import (
	"fmt"

	"hoangphuc.tech/hercules/domain/base"
	"hoangphuc.tech/hercules/infra/orm"
	"hoangphuc.tech/hercules/infra/orm/ext"
)

type ItemRepository struct {
	base.Repository
}

func (r *ItemRepository) FindAllOnPage(page uint, pageSize uint) ([]orm.Item, error) {
	var bqResult []map[string]interface{}
	result := DB().Table("v_product").Scopes(ext.Paginate(page, pageSize, "bq_created_at DESC, bq_id ASC, code ASC")).Find(&bqResult)

	var items []orm.Item
	for _, bqItem := range bqResult {
		i := orm.Item{
			Code:                fmt.Sprintf("%v", bqItem["code"]),
			Name:                fmt.Sprintf("%v", bqItem["name"]),
			MasterSKU:           fmt.Sprintf("%v", bqItem["master_sku"]),
			PrimaryCategoryCode: fmt.Sprintf("%v", bqItem["category_code"]),
			PrimaryCategory:     orm.Category{Code: fmt.Sprintf("%v", bqItem["category_code"])},
			BrandCode:           fmt.Sprintf("%v", bqItem["brand_code"]),
			Brand:               orm.Brand{Code: fmt.Sprintf("%v", bqItem["brand_code"])},
			Entity: orm.Entity{
				CreatedBy: "system",
				UpdatedBy: "system",
			},
		}

		i.VariantAttributes = new(ext.JSON)
		i.VariantAttributes.Scan(fmt.Sprintf("%v", bqItem["variant_attributes"]))

		items = append(items, i)
	}
	return items, result.Error
}

func (r *ItemRepository) Count() (uint64, error) {
	var c int64
	result := DB().Table("v_product").Count(&c)
	return uint64(c), result.Error
}
