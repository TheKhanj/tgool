package tgool

import (
	"path"

	"github.com/thekhanj/drouter"
)

type routeMetadata struct {
	path    string
	title   string
	hasBody bool

	controller Controller
	method     string
}

type RouteBuilder struct {
	controller  Controller
	prefixRoute string
	metadatas   []routeMetadata
}

func (r *RouteBuilder) setController(
	controller Controller,
) *RouteBuilder {
	r.controller = controller

	return r
}

func (r *RouteBuilder) SetPrefixRoute(
	prefixRoute string,
) *RouteBuilder {
	r.prefixRoute = prefixRoute

	return r
}

func (r *RouteBuilder) AddMethod(
	route string,
	methodName string,
) *RouteBuilder {
	path := path.Join(r.prefixRoute, route)
	r.metadatas = append(r.metadatas, routeMetadata{
		path:       path,
		method:     methodName,
		controller: r.controller,
		title:      "",
		hasBody:    false,
	})

	return r
}

func (r *RouteBuilder) WithTitle(title string) *RouteBuilder {
	r.metadatas[len(r.metadatas)-1].title = title

	return r
}

func (r *RouteBuilder) WithBody() *RouteBuilder {
	r.metadatas[len(r.metadatas)-1].hasBody = true

	return r
}

func (r *RouteBuilder) Build() *drouter.Router {
	router := drouter.New()

	for _, metadata := range r.metadatas {
		router.AddRoute(metadata.path, metadata)
	}

	return router
}
