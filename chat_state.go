package tgool

import (
	"path"

	"github.com/thekhanj/drouter"
)

type chatState struct {
	currentRoute string
}

func (c *chatState) GetPath() string {
	return c.currentRoute
}

func (c *chatState) SetPath(route string) *chatState {
	if route[0] != '/' {
		c.currentRoute = drouter.CleanPath(
			"/" + path.Join(
				c.GetPath(),
				route,
			),
		)
		return c
	}

	c.currentRoute = path.Join(route)

	return c
}

type chatsState struct {
	chats map[int64]*chatState
}

func (c *chatsState) GetChat(chatId int64) *chatState {
	if c.chats == nil {
		c.chats = make(map[int64]*chatState)
	}

	chat, ok := c.chats[chatId]
	if !ok {
		c.chats[chatId] = &chatState{
			currentRoute: "/",
		}
		chat = c.chats[chatId]
	}

	return chat
}
