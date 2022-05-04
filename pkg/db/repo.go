package db

type Repository[T any, K comparable] interface {
	Find(id K) (*T, bool)
	Create(model *T) error
	Update(model *T) error
	Delete(id K) error
}
