package base

type Repository interface {
	Validate(entity interface{}) error
	GetByID(id uint64) (interface{}, error)
	Search(query interface{}) ([]interface{}, error)
	Create(entity *interface{}) error
	Update(entity *interface{}) error
	Delete(id uint64) error
}
