package tgool

import (
	"errors"
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

func (r *RouteBuilder) SetController(
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

func (r *RouteBuilder) Build() (*drouter.Router, error) {
	router := drouter.New()

	for _, metadata := range r.metadatas {
		if metadata.controller == nil {
			return nil, errors.New("controller is not set")
		}

		router.AddRoute(metadata.path, metadata)
	}

	return router, nil
}
