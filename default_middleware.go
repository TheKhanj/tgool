package tgool

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type DefaultMiddleWare struct{}

var _ Middleware = &DefaultMiddleWare{}

func (m *DefaultMiddleWare) Handle(
	ctx Context,
	next func(),
) tg.Chattable {
	logrus.Info("[DefaultMiddleware] no route matched")

	return nil
}
