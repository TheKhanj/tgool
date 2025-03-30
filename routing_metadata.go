package tgool

import (
	"log"
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

type RouterBuilder struct {
	controller  Controller
	prefixRoute string
	metadatas   []routeMetadata
}

func (r *RouterBuilder) setController(
	controller Controller,
) *RouterBuilder {
	r.controller = controller

	return r
}

func (r *RouterBuilder) SetPrefixRoute(
	prefixRoute string,
) *RouterBuilder {
	r.prefixRoute = prefixRoute

	return r
}

func (r *RouterBuilder) AddMethod(
	route string,
	methodName string,
) *RouterBuilder {
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

func (r *RouterBuilder) WithTitle(title string) *RouterBuilder {
	r.metadatas[len(r.metadatas)-1].title = title

	return r
}

func (r *RouterBuilder) WithBody() *RouterBuilder {
	r.metadatas[len(r.metadatas)-1].hasBody = true

	return r
}

func (r *RouterBuilder) Build() *drouter.Router {
	router := drouter.New()

	for _, metadata := range r.metadatas {
		log.Printf("tgool: added route %s", metadata.path)
		router.AddRoute(metadata.path, metadata)
	}

	return router
}
