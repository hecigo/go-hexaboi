package model

type ItemAttribute struct {
	Code      string `json:"brand_code"`
	Name      string `json:"brand_name"`
	IsActived int8   `json:"is_actived"`
	IsDeleted int8   `json:"is_deleted"`
}

func (b ItemAttribute) String() string {
	return b.Code + "\t" + b.Name
}
