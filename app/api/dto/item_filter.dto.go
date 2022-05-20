package dto

type ItemFilterDto struct {
	Items     []string `query:"items" validate:"required"`
	PageIndex int      `query:"page_index" validate:"required"`
	PageSize  int      `query:"page_size" validate:"min=0,max=32"`
}
