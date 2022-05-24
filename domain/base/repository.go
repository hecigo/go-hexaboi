package base

type Repository interface {
	Validate(entity interface{}) error
	GetByID(id uint) (interface{}, error)
	Search(query interface{}) ([]interface{}, error)
	Create(entity interface{}) (interface{}, error)
	Update(entity interface{}) (interface{}, error)
	Delete(id uint) error
}
