package schemas

import "github.com/getkin/kin-openapi/openapi3"

var UserSchema = openapi3.NewSchema().
	WithProperties(map[string]*openapi3.Schema{
		"id":        openapi3.NewUUIDSchema(),
		"name":      openapi3.NewStringSchema().WithFormat("text").WithMin(3),
		"email":     openapi3.NewStringSchema().WithFormat("email").WithPattern("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"),
		"createdAt": openapi3.NewDateTimeSchema(),
		"updatedAt": openapi3.NewDateTimeSchema(),
	}).
	WithRequired([]string{
		"id",
		"name",
		"email",
		"createdAt",
		"updatedAt",
	})

var CreateUserSchema = openapi3.NewSchema().
	WithProperties(map[string]*openapi3.Schema{
		"name":  openapi3.NewStringSchema().WithFormat("text").WithMin(3),
		"email": openapi3.NewStringSchema().WithFormat("email").WithPattern("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"),
	}).
	WithRequired([]string{
		"name",
		"email",
	})

var UpdateUserSchema = openapi3.NewSchema().
	WithProperties(map[string]*openapi3.Schema{
		"name":  openapi3.NewStringSchema().WithFormat("text").WithMin(3),
		"email": openapi3.NewStringSchema().WithFormat("email").WithPattern("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"),
	}).
	WithRequired([]string{
		"name",
		"email",
	})
