package http

import (
	"fmt"
	"regexp"

	"github.com/connor-davis/dynamic-crud/cmd/api/http/routes"
	"github.com/connor-davis/dynamic-crud/common"
	"github.com/connor-davis/dynamic-crud/internal/routing"
	"github.com/connor-davis/dynamic-crud/internal/routing/schemas"
	"github.com/connor-davis/dynamic-crud/internal/storage"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

type HttpRouter interface {
	InitializeRoutes(router fiber.Router)
	InitializeOpenAPI() *openapi3.T
}

type httpRouter struct {
	storage storage.Storage
	routes  []routing.Route
}

func NewHttpRouter(storage storage.Storage) HttpRouter {
	usersRouter := routes.NewUsersRouter(storage)
	usersRoutes := usersRouter.LoadRoutes()

	routes := []routing.Route{}

	routes = append(routes, usersRoutes...)

	return &httpRouter{
		storage: storage,
		routes:  routes,
	}
}

func (h *httpRouter) InitializeRoutes(router fiber.Router) {
	for _, route := range h.routes {
		path := regexp.MustCompile(`\{([^}]+)\}`).ReplaceAllString(route.Path, ":$1")

		switch route.Method {
		case routing.GET:
			router.Get(path, append(route.Middlewares, route.Handler)...)
		case routing.POST:
			router.Post(path, append(route.Middlewares, route.Handler)...)
		case routing.PUT:
			router.Put(path, append(route.Middlewares, route.Handler)...)
		case routing.DELETE:
			router.Delete(path, append(route.Middlewares, route.Handler)...)
		}
	}
}

func (h *httpRouter) InitializeOpenAPI() *openapi3.T {
	paths := openapi3.NewPaths()

	schemas := openapi3.Schemas{
		"SuccessResponse": schemas.SuccessSchema.NewRef(),
		"ErrorResponse":   schemas.ErrorSchema.NewRef(),
		"User":            schemas.UserSchema.NewRef(),
	}

	for _, route := range h.routes {
		pathItem := &openapi3.PathItem{}

		switch route.Method {
		case routing.GET:
			pathItem.Get = &openapi3.Operation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        route.Tags,
				Parameters:  route.Parameters,
				Responses:   route.Responses,
			}
		case routing.POST:
			if route.CreateSchema == nil {
				continue
			}

			schemas[fmt.Sprintf("Create%s", route.Entity)] = route.CreateSchema.NewRef()

			pathItem.Post = &openapi3.Operation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        route.Tags,
				Parameters:  route.Parameters,
				RequestBody: route.RequestBody,
				Responses:   route.Responses,
			}
		case routing.PUT:
			if route.UpdateSchema == nil {
				continue
			}

			schemas[fmt.Sprintf("Update%s", route.Entity)] = route.UpdateSchema.NewRef()

			pathItem.Put = &openapi3.Operation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        route.Tags,
				Parameters:  route.Parameters,
				RequestBody: route.RequestBody,
				Responses:   route.Responses,
			}
		case routing.DELETE:
			pathItem.Delete = &openapi3.Operation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        route.Tags,
				Parameters:  route.Parameters,
				Responses:   route.Responses,
			}
		}

		path := fmt.Sprintf("/api%s", route.Path)

		existingPathItem := paths.Find(path)

		if existingPathItem != nil {
			switch route.Method {
			case routing.GET:
				existingPathItem.Get = pathItem.Get
			case routing.POST:
				existingPathItem.Post = pathItem.Post
			case routing.PUT:
				existingPathItem.Put = pathItem.Put
			case routing.DELETE:
				existingPathItem.Delete = pathItem.Delete
			}
		} else {
			paths.Set(path, pathItem)
		}
	}

	return &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:   common.EnvString("APP_NAME", "Dynamic CRUD API"),
			Version: common.EnvString("APP_VERSION", "1.0.0"),
		},
		Servers: openapi3.Servers{
			{
				URL:         fmt.Sprintf("http://localhost:%s", common.EnvString("APP_PORT", "6173")),
				Description: "Development",
			},
			{
				URL:         common.EnvString("APP_BASE_URL", "https://example.com"),
				Description: "Production",
			},
		},
		Tags:  openapi3.Tags{},
		Paths: paths,
		Components: &openapi3.Components{
			Schemas: schemas,
		},
	}
}
