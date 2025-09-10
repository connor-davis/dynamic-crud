package routes

import "github.com/connor-davis/dynamic-crud/internal/routing"

type Router interface {
	LoadRoutes() []routing.Route
}
