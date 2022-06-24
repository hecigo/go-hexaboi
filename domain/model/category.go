package model

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

// Kind of how to group products.
type Division string

const (
	DIVISION_MERCHANDISE Division = "merchandise"
	DIVISION_CONSUMER    Division = "consumer"
	DIVISION_CAMPAIGN    Division = "campaign"
	DIVISION_CUSTOM      Division = "custom"
)
