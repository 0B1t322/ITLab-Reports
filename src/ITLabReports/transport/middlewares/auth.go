package middlewares

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	mcontext "github.com/RTUITLab/ITLab-Reports/transport/middlewares/context"
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares/rolegetter"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitloggrus "github.com/go-kit/kit/log/logrus"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

var (
	NotAdmin = errors.New("Not admin")
	NotSuperAdmin = errors.New("Not super admin")
	TokenNotValid = errors.New("Token is not valid")
	TokenExpired = errors.New("Token expired")
	FailedToParseToken = errors.New("Failed to parse token")
	RoleNotFound = errors.New("Role is not found")
)

type authJWKS struct {
	UserRole       	string
	AdminRole      	string
	SuperAdminRole 	string
	RoleClaim      	string
	JWKSUrl			string

	jwks			*keyfunc.JWKS
	refreshTime		time.Duration

	auth			MiddlewareWithContext[any, any]
	admin			MiddlewareWithContext[any, any]
	superAdmin		MiddlewareWithContext[any, any]

	logger 			log.Logger

	roleGetter		*rolegetter.RoleGetter
}


// Check that token is valid
// 
// put into middleware context
// 
// catchable errors:
// 	TokenNotValid
// 	RoleNotFound
func (a *authJWKS) Auth() MiddlewareWithContext[any, any] {
	return a.auth
}

// Check that user have role admin
// 
// catchable errors:
// 	NotAdmin
func (a *authJWKS) IsAdmin() MiddlewareWithContext[any, any] {
	return a.admin
}

// Check that user have role superadmin
// 
// catchable errors:
// 	NotSuperAdmin
func (a *authJWKS) IsSuperAdmin() MiddlewareWithContext[any, any] {
	return a.superAdmin
}

func (a *authJWKS) setUserRole(userRole string) {
	a.UserRole = userRole
}

func (a *authJWKS) setAdminRole(admin string) {
	a.AdminRole = admin
}

func (a *authJWKS) setSuperAdminRole(superAdmin string) {
	a.SuperAdminRole = superAdmin
}

func (a *authJWKS) setJWKSUrl(url string) {
	a.JWKSUrl = url
}
func (a *authJWKS) setRoleClaim(claim string) {
	a.RoleClaim = claim
}

func (a *authJWKS) setRefreshTime(refreshTime time.Duration) {
	a.refreshTime = refreshTime
}

func (a *authJWKS) setLogger(logger log.Logger) {
	a.logger = logger
}

func NewJWKSAuth(
	opts ...AutherOptions,
) Auther {
	a := &authJWKS{
		UserRole: "user",
		AdminRole: "admin",
		SuperAdminRole: "superadmin",
		RoleClaim: "claim",
		JWKSUrl: "https://examaple.com",
		refreshTime: time.Minute,
		logger: kitloggrus.NewLogger(logrus.StandardLogger()),
	}

	for _, opt := range opts {
		opt(a)
	}

	a.logger = log.With(a.logger, "from", "AuthMiddleware")

	a.roleGetter = rolegetter.New(
		a.SuperAdminRole,
		a.AdminRole,
		a.UserRole,
	)

	if err := a.buildJWKS(); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"from": "NewJWKSAuth",
				"err": err,
			},
		).Panic("Failed to build jwks")
	}

	a.buildMiddlewares()
	return a
}

func (a *authJWKS) buildJWKS() error {
	jwks, err := keyfunc.Get(
		a.JWKSUrl,
		keyfunc.Options{
			RefreshInterval: a.refreshTime,
			RefreshErrorHandler: func(err error) {
				level.Error(a.logger).Log("err", err)
			},
			
		},
	)
	if err != nil {
		return err
	}

	a.jwks = jwks

	return nil
}

func (a *authJWKS) buildMiddlewares() {
	a.auth = func(next EndpointWithContext[any, any]) EndpointWithContext[any, any] {
		return func(ctx mcontext.MiddlewareContext, req any) (any, error) {
			_t, err := ctx.GetToken()
			if err != nil {
				return nil, TokenNotValid
			}

			jwtToken := strings.ReplaceAll(_t, "Bearer ", "")

			var claims jwt.MapClaims

			token, err := jwt.ParseWithClaims(jwtToken, &claims, a.jwks.Keyfunc)
			if validErr, ok := err.(*jwt.ValidationError); ok {
				switch validErr.Errors {
				case jwt.ValidationErrorExpired:
					return nil, TokenExpired
				default:
					level.Error(a.logger).Log("err", validErr.Error())
				}
			} else if err != nil {
				level.Error(a.logger).Log("err", err)
				return nil, FailedToParseToken
			}


			if !token.Valid {
				return nil, TokenNotValid
			}

			role, err := a.roleGetter.GetRole(claims, a.RoleClaim)
			if err != nil {
				return nil, RoleNotFound
			}
			ctx.SetRole(role)

			ctx.SetUserID(fmt.Sprint(claims["sub"]))

			return next(ctx, req)
		}
	}

	a.admin = func(next EndpointWithContext[any, any]) EndpointWithContext[any, any] {
		return func(ctx mcontext.MiddlewareContext, req any) (any, error) {
			role, err := ctx.GetRole()
			if err != nil {
				return nil, err
			}

			if !(role == a.AdminRole) {
				return nil, NotAdmin
			}

			return next(ctx, req)
		}
	}

	a.superAdmin = func(next EndpointWithContext[any, any]) EndpointWithContext[any, any] {
		return func(ctx mcontext.MiddlewareContext, req any) (any, error) {
			role, err := ctx.GetRole()
			if err != nil {
				return nil, err
			}

			if !(role == a.SuperAdminRole) {
				return nil, NotSuperAdmin
			}

			return next(ctx, req)
		}
	}
}

type testAuth struct {
	UserRole       	string
	AdminRole      	string
	SuperAdminRole 	string
	RoleClaim      	string

	auth			MiddlewareWithContext[any, any]
	admin			MiddlewareWithContext[any, any]
	superAdmin		MiddlewareWithContext[any, any]

	logger 			log.Logger

	roleGetter		*rolegetter.RoleGetter
}

func (a *testAuth) setUserRole(userRole string) {
	a.UserRole = userRole
}

func (a *testAuth) setAdminRole(admin string) {
	a.AdminRole = admin
}

func (a *testAuth) setSuperAdminRole(superAdmin string) {
	a.SuperAdminRole = superAdmin
}

func (a *testAuth) setJWKSUrl(url string) {
	
}
func (a *testAuth) setRoleClaim(claim string) {
	a.RoleClaim = claim
}

func (a *testAuth) setRefreshTime(refreshTime time.Duration) {
	
}

func (a *testAuth) setLogger(logger log.Logger) {
	a.logger = logger
}

// Check that token is valid
// 
// put into middleware context
// 
// catchable errors:
// 	TokenNotValid
// 	RoleNotFound
func (a *testAuth) Auth() MiddlewareWithContext[any, any] {
	return a.auth
}

// Check that user have role admin
// 
// catchable errors:
// 	NotAdmin
func (a *testAuth) IsAdmin() MiddlewareWithContext[any, any] {
	return a.admin
}

// Check that user have role superadmin
// 
// catchable errors:
// 	NotSuperAdmin
func (a *testAuth) IsSuperAdmin() MiddlewareWithContext[any, any] {
	return a.superAdmin
}

func NewTestAuth(
	opts ...AutherOptions,
) Auther {
	a := &testAuth{
		UserRole: "user",
		AdminRole: "admin",
		SuperAdminRole: "superadmin",
		RoleClaim: "claim",
		logger: kitloggrus.NewLogger(logrus.StandardLogger()),
	}

	for _, opt := range opts {
		opt(a)
	}

	a.logger = log.With(a.logger, "from", "AuthMiddleware")
	logrus.Info("auth cfg ", a)
	a.roleGetter = rolegetter.New(
		a.SuperAdminRole,
		a.AdminRole,
		a.UserRole,
	)

	a.buildMiddlewares()
	return a
}

func (a *testAuth) buildMiddlewares() {
	a.auth = func(next EndpointWithContext[any, any]) EndpointWithContext[any, any] {
		return func(ctx mcontext.MiddlewareContext, req any) (any, error) {
			_t, err := ctx.GetToken()
			if err != nil {
				return nil, TokenNotValid
			}

			jwtToken := strings.ReplaceAll(_t, "Bearer ", "")
			
			var claims jwt.MapClaims

			_, err = jwt.ParseWithClaims(
				jwtToken, 
				&claims,
				func(t *jwt.Token) (interface{}, error) {
					if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
					}
					return []byte("test"), nil
				},
			)
			if validErr, ok := err.(*jwt.ValidationError); ok {
				switch validErr.Errors {
				case jwt.ValidationErrorExpired:
					// Pass
					// return nil, TokenExpired
				default:
					level.Error(a.logger).Log("err", validErr.Error())
				}
			} else if err != nil {
				level.Error(a.logger).Log("err", err)
				return nil, FailedToParseToken
			}


			// if !token.Valid {
			// 	return nil, TokenNotValid
			// }

			role, err := a.roleGetter.GetRole(claims, a.RoleClaim)
			if err != nil {
				return nil, RoleNotFound
			}
			ctx.SetRole(role)

			ctx.SetUserID(fmt.Sprint(claims["sub"]))

			return next(ctx, req)
		}
	}

	a.admin = func(next EndpointWithContext[any, any]) EndpointWithContext[any, any] {
		return func(ctx mcontext.MiddlewareContext, req any) (any, error) {
			role, err := ctx.GetRole()
			if err != nil {
				return nil, err
			}

			if !(role == a.AdminRole) {
				return nil, NotAdmin
			}

			return next(ctx, req)
		}
	}

	a.superAdmin = func(next EndpointWithContext[any, any]) EndpointWithContext[any, any] {
		return func(ctx mcontext.MiddlewareContext, req any) (any, error) {
			role, err := ctx.GetRole()
			if err != nil {
				return nil, err
			}

			if !(role == a.SuperAdminRole) {
				return nil, NotSuperAdmin
			}

			return next(ctx, req)
		}
	}
}