package ext

import (
	"gorm.io/gorm"
)

func Paginate(page uint, pageSize uint, order interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		if pageSize == 0 {
			pageSize = 1
		}

		offset := int((page - 1) * pageSize)
		return db.Offset(offset).Limit(int(pageSize)).Order(order)
	}
}
