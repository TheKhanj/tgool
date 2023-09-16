package tgool

import (
	"fmt"
	"reflect"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/thekhanj/drouter"
)

type ControllerMiddleware struct {
	router     *drouter.Router
}

var _ Middleware = &ControllerMiddleware{}

func NewControllerMiddleware(
	controllers ...Controller,
) (*ControllerMiddleware, error) {
	builder := RouteBuilder{
		prefixRoute: "/",
	}

	for _, controller := range controllers {
		controller.AddRoutes(&builder)
	}

	router, err := builder.Build()
	if err != nil {
		return nil, err
	}

	return &ControllerMiddleware{
		router: router,
	}, nil
}

func (m *ControllerMiddleware) Handle(
	ctx Context,
	next func(),
) tg.Chattable {
	path, method, methodFound := m.getHandler(ctx)

	if !methodFound {
		next()
		return nil
	}

	ctx.ChatsState().GetChat(ctx.GetChatId()).SetPath(path)

	return m.callMethod(ctx, method)
}

func (m *ControllerMiddleware) callMethod(
	ctx Context,
	method reflect.Value,
) tg.Chattable {
	chatId := ctx.GetChatId()

	passingContext := reflect.ValueOf(ctx)
	args := []reflect.Value{
		passingContext,
	}

	ret := method.Call(args)
	result := ret[0].Interface()
	resultError := ret[1].Interface()

	if err, ok := resultError.(error); ok {
		return tg.NewMessage(chatId, err.Error())
	}

	if chattable, ok := result.(tg.Chattable); ok {
		return chattable
	}

	logrus.Debug("result of controller is not being handled")
	return nil
}

func (m *ControllerMiddleware) getHandler(
	ctx Context,
) (path string, method reflect.Value, ok bool) {
	chatId := ctx.GetChatId()
	chatState := ctx.ChatsState().GetChat(chatId)

	ok = false

	params := make(drouter.Params, 0, 0)
	currentPath := chatState.GetPath()
	logrus.Warn("current path: ", currentPath)

	handle, _ := m.router.Lookup(currentPath, &params)
	if handle != nil {
		route := handle.(routeMetadata)
		if route.hasBody == true {
			logrus.Info(route.path)
			method, ok = getMethodByRouteMetadata(&route)
			return route.path, method, ok
		}
	}

	handle, _ = m.router.Lookup(ctx.GetRoute(), &params)
	if handle != nil {
		route := handle.(routeMetadata)
		logrus.Info(route.path)
		method, ok = getMethodByRouteMetadata(&route)
		return route.path, method, ok
	}

	return "", method, false
}

func getMethodByRouteMetadata(route *routeMetadata) (method reflect.Value, ok bool) {
	method = reflect.
		ValueOf(route.controller).
		MethodByName(route.method)

	if !method.IsValid() {
		logrus.Debug(
			fmt.Sprintf("method %s is invalid", route.method),
		)
		return method, false
	}
	return method, true
}
