package orm

import "hoangphuc.tech/go-hexaboi/domain/model"

// Item grouping, category
type Category struct {
	Code      string `json:"category_code"`
	Name      string `json:"category_name"`
	IsActived int8   `json:"is_actived"`
	IsDeleted int8   `json:"is_deleted"`
}

func (b Category) String() string {
	return b.Code + "\t" + b.Name
}

func (c *Category) ToModel(cate *model.Category) {
	if cate == nil {
		cate = new(model.Category)
	}
	cate.Code = c.Code
	cate.Name = c.Name
	cate.IsActived = c.IsActived
	cate.IsDeleted = c.IsDeleted
}
