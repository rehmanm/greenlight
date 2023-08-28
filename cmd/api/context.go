package main

import (
	"context"
	"net/http"

	"github.com/rehmanm/greenlight/internal/data"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	context := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(context)
}

func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in context")
	}
	return user
}
