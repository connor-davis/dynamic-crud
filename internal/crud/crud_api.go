package crud

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/connor-davis/dynamic-crud/internal/routing"
	"github.com/connor-davis/dynamic-crud/internal/routing/schemas"
	"github.com/connor-davis/dynamic-crud/internal/storage"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CrudApi[T any] interface {
	AssignCreateSchema(schema *openapi3.Schema) CrudApi[T]
	AssignUpdateSchema(schema *openapi3.Schema) CrudApi[T]
	CreateRoute() routing.Route
	UpdateRoute() routing.Route
	DeleteRoute() routing.Route
	GetOneRoute() routing.Route
	GetAllRoute() routing.Route
}

type crudApi[T any] struct {
	storage storage.Storage
	name    string
	crud    Crud[T]
	create  *openapi3.Schema
	update  *openapi3.Schema
}

type UpdateParams struct {
	Id string `json:"id"`
}

type DeleteParams struct {
	Id string `json:"id"`
}

type GetOneParams struct {
	Id string `json:"id"`
}

func NewCrudApi[T any](storage storage.Storage) CrudApi[T] {
	crud := NewCrud[T](storage)

	tReflection := reflect.TypeOf(new(T))
	tReflectionName := tReflection.Elem().Name()

	log.Printf("Initialized CRUD API for %s at /%ss", tReflectionName, strings.ToLower(tReflectionName))

	return &crudApi[T]{
		storage: storage,
		name:    tReflectionName,
		crud:    crud,
	}
}

func (c *crudApi[T]) AssignCreateSchema(schema *openapi3.Schema) CrudApi[T] {
	c.create = schema

	return c
}

func (c *crudApi[T]) AssignUpdateSchema(schema *openapi3.Schema) CrudApi[T] {
	c.update = schema

	return c
}

func (c *crudApi[T]) CreateRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("%s created successfully.", c.name)).
			WithContent(openapi3.Content{
				"text/plain": openapi3.NewMediaType().
					WithSchema(openapi3.NewStringSchema().WithDefault("OK")),
			}),
	})

	responses.Set("400", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorSchema).
			WithDescription("Bad Request").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorSchema),
			}),
	})

	responses.Set("401", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorSchema).
			WithDescription("Unauthorized").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorSchema),
			}),
	})

	responses.Set("403", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorSchema).
			WithDescription("Forbidden").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorSchema),
			}),
	})

	responses.Set("500", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorSchema).
			WithDescription("Internal Server Error").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorSchema),
			}),
	})

	return routing.Route{
		OpenAPIMetadata: routing.OpenAPIMetadata{
			Summary:     fmt.Sprintf("Create %s", c.name),
			Description: fmt.Sprintf("This endpoint creates a new %s.", strings.ToLower(c.name)),
			Tags:        []string{fmt.Sprintf("%ss", c.name)},
			Parameters:  nil,
			RequestBody: &openapi3.RequestBodyRef{
				Value: openapi3.NewRequestBody().
					WithRequired(true).
					WithJSONSchemaRef(c.create.NewRef()).
					WithDescription(fmt.Sprintf("Payload to create a new %s.", strings.ToLower(c.name))),
			},
			Responses: responses,
		},
		Entity:       c.name,
		CreateSchema: c.create,
		UpdateSchema: nil,
		Method:       routing.POST,
		Path:         fmt.Sprintf("/%ss", strings.ToLower(c.name)),
		Middlewares:  []fiber.Handler{},
		Handler: func(ctx *fiber.Ctx) error {
			var entity T

			if err := ctx.BodyParser(&entity); err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Bad Request",
					"message": err.Error(),
				})
			}

			if err := c.crud.Create(&entity); err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}

			return ctx.Status(fiber.StatusOK).SendString("OK")
		},
	}
}

func (c *crudApi[T]) UpdateRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("%s updated successfully", c.name)).
			WithContent(openapi3.Content{
				"text/plain": openapi3.NewMediaType().
					WithSchema(openapi3.NewStringSchema().WithDefault("OK")),
			}),
	})

	responses.Set("400", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorSchema).
			WithDescription("Bad Request").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorSchema),
			}),
	})

	responses.Set("401", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorSchema).
			WithDescription("Unauthorized").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorSchema),
			}),
	})

	responses.Set("403", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorSchema).
			WithDescription("Forbidden").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorSchema),
			}),
	})

	responses.Set("404", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorSchema).
			WithDescription("Not Found").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorSchema),
			}),
	})

	responses.Set("500", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorSchema).
			WithDescription("Internal Server Error").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorSchema),
			}),
	})

	return routing.Route{
		OpenAPIMetadata: routing.OpenAPIMetadata{
			Summary:     fmt.Sprintf("Update %s", c.name),
			Description: fmt.Sprintf("This endpoint updates an existing %s.", strings.ToLower(c.name)),
			Tags:        []string{fmt.Sprintf("%ss", c.name)},
			Parameters: []*openapi3.ParameterRef{
				{
					Value: openapi3.NewPathParameter("id").
						WithRequired(true).
						WithSchema(openapi3.NewUUIDSchema()),
				},
			},
			RequestBody: &openapi3.RequestBodyRef{
				Value: openapi3.NewRequestBody().
					WithRequired(true).
					WithJSONSchemaRef(c.update.NewRef()).
					WithDescription(fmt.Sprintf("Payload to update an existing %s.", strings.ToLower(c.name))),
			},
			Responses: responses,
		},
		Entity:       c.name,
		CreateSchema: nil,
		UpdateSchema: c.update,
		Method:       routing.PUT,
		Path:         fmt.Sprintf("/%ss/:id", strings.ToLower(c.name)),
		Middlewares:  []fiber.Handler{},
		Handler: func(ctx *fiber.Ctx) error {
			var params UpdateParams

			if err := ctx.ParamsParser(&params); err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Bad Request",
					"message": err.Error(),
				})
			}

			var entity T

			if err := ctx.BodyParser(&entity); err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Bad Request",
					"message": err.Error(),
				})
			}

			if err := c.crud.Update(params.Id, &entity); err != nil {
				if err == gorm.ErrRecordNotFound {
					return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
						"error":   "Not Found",
						"message": fmt.Sprintf("The %s was not found.", strings.ToLower(c.name)),
					})
				}

				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}

			return ctx.Status(fiber.StatusOK).SendString("OK")
		},
	}
}

func (c *crudApi[T]) DeleteRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("%s deleted successfully.", c.name)).
			WithContent(openapi3.Content{
				"text/plain": openapi3.NewMediaType().
					WithSchema(openapi3.NewStringSchema().WithDefault("OK")),
			}),
	})

	responses.Set("400", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorSchema).
			WithDescription("Bad Request").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorSchema),
			}),
	})

	responses.Set("401", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorSchema).
			WithDescription("Unauthorized").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorSchema),
			}),
	})

	responses.Set("403", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorSchema).
			WithDescription("Forbidden").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorSchema),
			}),
	})

	responses.Set("404", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorSchema).
			WithDescription("Not Found").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorSchema),
			}),
	})

	responses.Set("500", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorSchema).
			WithDescription("Internal Server Error").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorSchema),
			}),
	})

	return routing.Route{
		OpenAPIMetadata: routing.OpenAPIMetadata{
			Summary:     fmt.Sprintf("Delete %s", c.name),
			Description: fmt.Sprintf("This endpoint deletes an existing %s.", strings.ToLower(c.name)),
			Tags:        []string{fmt.Sprintf("%ss", c.name)},
			Parameters: []*openapi3.ParameterRef{
				{
					Value: openapi3.NewPathParameter("id").
						WithRequired(true).
						WithSchema(openapi3.NewUUIDSchema()),
				},
			},
			RequestBody: nil,
			Responses:   responses,
		},
		Entity:       c.name,
		CreateSchema: nil,
		UpdateSchema: nil,
		Method:       routing.DELETE,
		Path:         fmt.Sprintf("/%ss/:id", strings.ToLower(c.name)),
		Middlewares:  []fiber.Handler{},
		Handler: func(ctx *fiber.Ctx) error {
			var params DeleteParams

			if err := ctx.ParamsParser(&params); err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Bad Request",
					"message": err.Error(),
				})
			}

			if err := c.crud.Delete(params.Id, new(T)); err != nil {
				if err == gorm.ErrRecordNotFound {
					return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
						"error":   "Not Found",
						"message": fmt.Sprintf("The %s was not found.", strings.ToLower(c.name)),
					})
				}

				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}

			return ctx.Status(fiber.StatusOK).SendString("OK")
		},
	}
}

func (c *crudApi[T]) GetOneRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("%s retrieved successfully.", c.name)).
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.SuccessSchema),
			}),
	})

	responses.Set("400", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorSchema).
			WithDescription("Bad Request").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorSchema),
			}),
	})

	responses.Set("401", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorSchema).
			WithDescription("Unauthorized").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorSchema),
			}),
	})

	responses.Set("403", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorSchema).
			WithDescription("Forbidden").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorSchema),
			}),
	})

	responses.Set("404", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorSchema).
			WithDescription("Not Found").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorSchema),
			}),
	})

	responses.Set("500", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorSchema).
			WithDescription("Internal Server Error").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorSchema),
			}),
	})

	return routing.Route{
		OpenAPIMetadata: routing.OpenAPIMetadata{
			Summary:     fmt.Sprintf("Get %s", c.name),
			Description: fmt.Sprintf("This endpoint retrieves an existing %s.", strings.ToLower(c.name)),
			Tags:        []string{fmt.Sprintf("%ss", c.name)},
			Parameters: []*openapi3.ParameterRef{
				{
					Value: openapi3.NewPathParameter("id").
						WithRequired(true).
						WithSchema(openapi3.NewUUIDSchema()),
				},
			},
			RequestBody: nil,
			Responses:   responses,
		},
		Entity:       c.name,
		CreateSchema: nil,
		UpdateSchema: nil,
		Method:       routing.GET,
		Path:         fmt.Sprintf("/%ss/:id", strings.ToLower(c.name)),
		Middlewares:  []fiber.Handler{},
		Handler: func(ctx *fiber.Ctx) error {
			var params GetOneParams

			if err := ctx.ParamsParser(&params); err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Bad Request",
					"message": err.Error(),
				})
			}

			var entity T

			if err := c.crud.FindOne(params.Id, &entity); err != nil {
				if err == gorm.ErrRecordNotFound {
					return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
						"error":   "Not Found",
						"message": fmt.Sprintf("The %s was not found.", strings.ToLower(c.name)),
					})
				}

				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}

			return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
				"item": entity,
			})
		},
	}
}

func (c *crudApi[T]) GetAllRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("%s's retrieved successfully.", c.name)).
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.SuccessSchema),
			}),
	})

	responses.Set("400", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorSchema).
			WithDescription("Bad Request").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorSchema),
			}),
	})

	responses.Set("401", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorSchema).
			WithDescription("Unauthorized").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorSchema),
			}),
	})

	responses.Set("403", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorSchema).
			WithDescription("Forbidden").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorSchema),
			}),
	})

	responses.Set("500", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(schemas.ErrorSchema).
			WithDescription("Internal Server Error").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchema(schemas.ErrorSchema),
			}),
	})

	return routing.Route{
		OpenAPIMetadata: routing.OpenAPIMetadata{
			Summary:     fmt.Sprintf("Get %ss", c.name),
			Description: fmt.Sprintf("This endpoint retrieves a list of %ss.", strings.ToLower(c.name)),
			Tags:        []string{fmt.Sprintf("%ss", c.name)},
			Parameters:  nil,
			RequestBody: nil,
			Responses:   responses,
		},
		Entity:       c.name,
		CreateSchema: nil,
		UpdateSchema: nil,
		Method:       routing.GET,
		Path:         fmt.Sprintf("/%ss", strings.ToLower(c.name)),
		Middlewares:  []fiber.Handler{},
		Handler: func(ctx *fiber.Ctx) error {
			var entities []T

			if err := c.crud.FindAll(&entities); err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}

			return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
				"items": entities,
			})
		},
	}
}
