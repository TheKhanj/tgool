package tgool

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thekhanj/drouter"
)

type Context interface {
	Bot() *tg.BotAPI
	Update() *tg.Update
	Redirect(path string)
	GetMessage() *tg.Message
	GetRoute() string
	GetChatId() int64
	GetFrom() *tg.User
	GetTelegramUserId() int64
	Params() *drouter.Params
	// this is shit
	ChatsState() *chatsState
}

type context struct {
	bot        *tg.BotAPI
	u          *tg.Update
	chatsState *chatsState
	params     *drouter.Params
}

func (c *context) Bot() *tg.BotAPI {
	return c.bot
}

func (c *context) Update() *tg.Update {
	return c.u
}

func (c *context) Redirect(path string) {
	c.chatsState.GetChat(c.GetChatId()).SetPath(path)
}

func (c *context) GetMessage() *tg.Message {
	if c.u.Message != nil {
		return c.u.Message
	}
	if c.u.CallbackQuery != nil &&
		c.u.CallbackQuery.Message != nil {
		return c.u.CallbackQuery.Message
	}
	return nil
}

func (c *context) GetRoute() string {
	if c.Update().CallbackQuery != nil {
		return c.Update().CallbackQuery.Data
	}
	if c.Update().Message != nil {
		return c.Update().Message.Text
	}

	return ""
}

func (c *context) GetChatId() int64 {
	message := c.GetMessage()
	if message == nil {
		return 0
	}

	if message.Chat != nil {
		return message.Chat.ID
	}

	return 0
}

func (c *context) GetFrom() *tg.User {
	if c.u.Message != nil && c.u.Message.From != nil {
		return c.u.Message.From
	}
	if c.u.CallbackQuery != nil && c.u.CallbackQuery.From != nil {
		return c.u.CallbackQuery.From
	}
	return nil
}

func (c *context) GetTelegramUserId() int64 {
	from := c.GetFrom()

	if from == nil {
		return 0
	}

	return from.ID
}

func (c *context) ChatsState() *chatsState {
	return c.chatsState
}

func (c *context) Params() *drouter.Params {
	return c.params
}

var _ Context = &context{}
