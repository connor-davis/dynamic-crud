package routing

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

type OpenAPIMetadata struct {
	Summary     string
	Description string
	Tags        []string
	Parameters  []*openapi3.ParameterRef
	RequestBody *openapi3.RequestBodyRef
	Responses   *openapi3.Responses
}

type RouteMethod string

const (
	GET    RouteMethod = "GET"
	POST   RouteMethod = "POST"
	PUT    RouteMethod = "PUT"
	PATCH  RouteMethod = "PATCH"
	DELETE RouteMethod = "DELETE"
)

type Route struct {
	OpenAPIMetadata

	Entity string

	CreateSchema *openapi3.Schema
	UpdateSchema *openapi3.Schema

	Method      RouteMethod
	Path        string
	Middlewares []fiber.Handler
	Handler     fiber.Handler
}
