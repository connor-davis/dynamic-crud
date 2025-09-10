package crud

import "github.com/connor-davis/dynamic-crud/internal/storage"

type Crud[T any] interface {
	Create(entity *T) error
	Update(entityId any, entity *T) error
	Delete(entityId any, entity *T) error
	FindOne(entityId any, entity *T) error
	FindAll(entities *[]T) error
}

type crud[T any] struct {
	storage storage.Storage
}

func NewCrud[T any](storage storage.Storage) Crud[T] {
	return &crud[T]{
		storage: storage,
	}
}

func (c *crud[T]) Create(entity *T) error {
	return c.storage.Database().Create(entity).Error
}

func (c *crud[T]) Update(entityId any, entity *T) error {
	return c.storage.Database().Model(entity).Where("id = ?", entityId).Updates(entity).Error
}

func (c *crud[T]) Delete(entityId any, entity *T) error {
	return c.storage.Database().Where("id = ?", entityId).Delete(entity).Error
}

func (c *crud[T]) FindOne(entityId any, entity *T) error {
	return c.storage.Database().First(entity, "id = ?", entityId).Error
}

func (c *crud[T]) FindAll(entities *[]T) error {
	return c.storage.Database().Find(&entities).Error
}
