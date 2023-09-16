package tgool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	chatsState := chatsState{}
	chatState := chatsState.GetChat(0)

	tests := [][2]string{
		{"", "/"},
		{"/test", "/test"},
		{"./foo", "/test/foo"},
		{"../foo", "/test/foo"},
		{"../../foo", "/foo"},
		{"../..", "/"},
		{"..", "/"},
		{"foo", "/foo"},
		{"goo", "/foo/goo"},
		{"///test", "/test"},
		{"foo///", "/test/foo"},
		{"/test////foo///goo", "/test/foo/goo"},
	}

	for _, test := range tests {
		set := test[0]
		if set != "" {
			chatState.SetPath(set)
		}
		expected := test[1]
		actual := chatState.GetPath()
		assert.Equal(t, expected, actual)
	}
}
