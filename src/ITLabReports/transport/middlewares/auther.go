package middlewares

import (
	"time"

	"github.com/go-kit/kit/log"
)



type Auther interface {

	// Check that token is valid
	// 
	// put into middleware context
	// 
	// catchable errors:
	// 	TokenNotValid
	// 	RoleNotFound
	Auth() MiddlewareWithContext[any, any]

	// Check that user have role admin
	// 
	// catchable errors:
	// 	NotAdmin
	IsAdmin() MiddlewareWithContext[any, any]

	// Check that user have role superadmin
	// 
	// catchable errors:
	// 	NotSuperAdmin
	IsSuperAdmin() MiddlewareWithContext[any, any]
}

type auther interface {
	setUserRole(string)

	setAdminRole(string)

	setSuperAdminRole(string)

	setJWKSUrl(string)

	setRoleClaim(string)

	setRefreshTime(time.Duration)
	
	setLogger(log.Logger)
}

type AutherOptions func(auther)

func WithUserRole(user string) AutherOptions {
	return func(a auther) {
		a.setUserRole(user)
	}
}

func WithAdminRole(admin string) AutherOptions {
	return func(a auther) {
		a.setAdminRole(admin)
	}
}

func WithSuperAdminRole(superAdmin string) AutherOptions {
	return func(a auther) {
		a.setSuperAdminRole(superAdmin)
	}
}

func WithJWKSUrl(url string) AutherOptions {
	return func(a auther) {
		a.setJWKSUrl(url)
	}
}

func WithRoleClaim(claim string) AutherOptions {
	return func(a auther) {
		a.setRoleClaim(claim)
	}
}

func WithRefreshTime(refreshTime time.Duration) AutherOptions {
	return func(a auther) {
		a.setRefreshTime(refreshTime)
	}
}

func WithLogger(logger log.Logger) AutherOptions {
	return func(a auther) {
		a.setLogger(logger)
	}
}