package routes

import (
	"github.com/connor-davis/dynamic-crud/internal/crud"
	"github.com/connor-davis/dynamic-crud/internal/models"
	"github.com/connor-davis/dynamic-crud/internal/routing"
	"github.com/connor-davis/dynamic-crud/internal/routing/schemas"
	"github.com/connor-davis/dynamic-crud/internal/storage"
)

type UsersRouter struct {
	storage storage.Storage
}

func NewUsersRouter(storage storage.Storage) Router {
	return &UsersRouter{
		storage: storage,
	}
}

func (r *UsersRouter) LoadRoutes() []routing.Route {
	crudApi := crud.NewCrudApi[models.User](r.storage).
		AssignCreateSchema(schemas.UserSchema).
		AssignUpdateSchema(schemas.UserSchema)

	getAllRoute := crudApi.GetAllRoute()
	getOneRoute := crudApi.GetOneRoute()
	createRoute := crudApi.CreateRoute()
	updateRoute := crudApi.UpdateRoute()
	deleteRoute := crudApi.DeleteRoute()

	return []routing.Route{
		getAllRoute,
		getOneRoute,
		createRoute,
		updateRoute,
		deleteRoute,
	}
}
