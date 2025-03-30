package tgool

type Controller interface {
	AddRoutes(builder *RouterBuilder)
}
