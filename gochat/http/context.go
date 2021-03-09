package http

import (
	"context"

	"github.com/0xdod/gochat/gochat"
)

func NewContextWithSession(ctx context.Context, s map[string]interface{}) context.Context {
	return context.WithValue(ctx, "session", s)
}

func SessionFromContext(ctx context.Context) map[interface{}]interface{} {
	session, ok := ctx.Value("session").(map[interface{}]interface{})
	if !ok {
		return nil
	}
	return session
}

func FlashFromContext(ctx context.Context) []FlashMessage {
	messages, ok := ctx.Value("messages").([]FlashMessage)
	if !ok {
		return nil
	}

	return messages
}

func NewContextWithUser(ctx context.Context, user *gochat.User) context.Context {
	return context.WithValue(ctx, "user", user)
}

func UserFromContext(ctx context.Context) *gochat.User {
	user, ok := ctx.Value("user").(*gochat.User)
	if !ok {
		return nil
	}
	return user
}
