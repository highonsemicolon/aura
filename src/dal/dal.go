package dal

type DAL[T any] interface {
	GetByID(id string) (*T, error)
	GetAll() ([]T, error)
	Create(entity *T) error
	Update(entity *T) error
	Delete(id string) error
}
