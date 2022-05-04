package db

type Repository[T any, K comparable] interface {
	FindOne(id K) (*T, bool)
	GetAll(qty, page int) []*T
	Create(model *T) error
	Update(model *T) error
	Delete(id K) error
}
