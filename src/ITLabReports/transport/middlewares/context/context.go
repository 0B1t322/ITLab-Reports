package context

import (
	"context"
	"fmt"
)

type MiddlewareContext interface {
	context.Context
	UsererContext
	TokenerContext
}

type UsererContext interface {
	// Return user id if have
	// if not have return error
	GetUserID() (string, error)

	// Set user id
	SetUserID(string)

	SetRole(string)

	GetRole() (string, error)
}

type TokenerContext interface {
	// Return token if have
	// if not return error
	GetTokenContext

	// set token
	SetToken(string)
}

type GetTokenContext interface {
	GetToken() (string, error)
}

type middlewareContext struct {
	context.Context

	userId string
	token  string
	role string
}

func GetFrom(ctx context.Context) (MiddlewareContext, error) {
	mctx, ok := ctx.(MiddlewareContext)
	if !ok {
		return nil, fmt.Errorf("is not middleware context")
	}

	return mctx, nil
}

func CreateOrGetFrom(ctx context.Context) (mctx MiddlewareContext) {
	var err error

	mctx, err = GetFrom(ctx)
	if err != nil {
		return New(ctx)
	}

	return mctx
}

func New(ctx context.Context) MiddlewareContext {
	return &middlewareContext{
		Context: ctx,
	}
}

func (m *middlewareContext) GetUserID() (string, error) {
	if m.userId == "" {
		return "", fmt.Errorf("UserID not set")
	}

	return m.userId, nil
}

func (m *middlewareContext) SetUserID(id string) {
	m.userId = id
}

func (m *middlewareContext) GetToken() (string, error) {
	if m.token == "" {
		return "", fmt.Errorf("Token not set")
	}
	return m.token, nil
}

func (m *middlewareContext) SetToken(token string) {
	m.token = token
}


func (m *middlewareContext) SetRole(role string) {
	m.role = role
}

func (m *middlewareContext) GetRole() (string, error) {
	if m.role == "" {
		return "", fmt.Errorf("Role is not set")
	}
	return m.role, nil
}