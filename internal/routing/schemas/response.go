package schemas

import (
	"github.com/getkin/kin-openapi/openapi3"
)

var ErrorSchema = openapi3.NewSchema().
	WithProperties(map[string]*openapi3.Schema{
		"error":   openapi3.NewStringSchema().WithFormat("text"),
		"message": openapi3.NewStringSchema().WithFormat("text"),
	}).
	WithRequired([]string{
		"error",
		"message",
	})

var SuccessSchema = openapi3.NewSchema().
	WithProperties(map[string]*openapi3.Schema{
		"item": openapi3.NewAnyOfSchema(
			UserSchema,
		),
		"items": openapi3.NewArraySchema().WithItems(
			openapi3.NewAnyOfSchema(
				UserSchema,
			),
		),
	}).
	WithRequired([]string{})
