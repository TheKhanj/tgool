package tgool

import (
	"fmt"
	"log"
	"reflect"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thekhanj/drouter"
)

type ControllerMiddleware struct {
	router *drouter.Router
}

var _ Middleware = &ControllerMiddleware{}

func NewControllerMiddleware(
	controllers ...Controller,
) *ControllerMiddleware {
	routeBuilder := RouterBuilder{
		prefixRoute: "/",
	}

	for _, controller := range controllers {
		controller.AddRoutes(routeBuilder.setController(controller))
	}

	router := routeBuilder.Build()

	return &ControllerMiddleware{router}
}

func (this *ControllerMiddleware) Handle(
	ctx Context,
	next func(),
) tg.Chattable {
	path, method, methodFound := this.getHandler(ctx)

	if !methodFound {
		next()
		return nil
	}

	ctx.ChatsState().GetChat(ctx.GetChatId()).SetPath(path)

	return this.callMethod(ctx, method)
}

func (this *ControllerMiddleware) callMethod(
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

	log.Println("tgool: result of controller is not being handled")
	return nil
}

func (this *ControllerMiddleware) getHandler(
	ctx Context,
) (path string, method reflect.Value, ok bool) {
	chatId := ctx.GetChatId()
	chatState := ctx.ChatsState().GetChat(chatId)

	ok = false

	// TODO: this is simply wrong! fix it!!!!
	p := make(drouter.Params, 0, 20)
	// TODO: this is shit
	ctx.(*context).params = &p
	currentPath := chatState.GetPath()
	log.Printf("tgool: current path is %s", currentPath)

	handle, _ := this.router.Lookup(currentPath, &p)
	if handle != nil {
		route := handle.(routeMetadata)
		if route.hasBody == true {
			log.Println("tgool: incoming request", currentPath)
			method, ok = getMethodByRouteMetadata(&route)
			return currentPath, method, ok
		}
	}

	p = make(drouter.Params, 0, 20)
	ctx.(*context).params = &p
	handle, _ = this.router.Lookup(ctx.GetRoute(), &p)
	if handle != nil {
		route := handle.(routeMetadata)
		log.Println("tgool: incoming request", ctx.GetRoute())
		method, ok = getMethodByRouteMetadata(&route)
		return ctx.GetRoute(), method, ok
	}

	return "", method, false
}

func getMethodByRouteMetadata(route *routeMetadata) (method reflect.Value, ok bool) {
	method = reflect.
		ValueOf(route.controller).
		MethodByName(route.method)

	if !method.IsValid() {
		log.Println(
			fmt.Sprintf("tgool: method %s is invalid", route.method),
		)
		return method, false
	}
	return method, true
}
