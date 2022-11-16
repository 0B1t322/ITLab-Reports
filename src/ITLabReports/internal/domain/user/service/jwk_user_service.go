package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/RTUITLab/ITLab-Reports/internal/common/errors"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/RTUITLab/ITLab-Reports/internal/models/entity"
	"github.com/golang-jwt/jwt/v4"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

type JWKAuthServiceConfig struct {
	UserRole       string
	AdminRole      string
	SuperAdminRole string

	RoleClaim   string
	ScopeClaim  string
	UserIDClaim string

	KeyURL string

	RequiredScope string

	RefreshTime time.Duration
}

type JWKUserService struct {
	jwks *keyfunc.JWKS

	RoleClaim  string
	roleGetter *RoleGetter

	ScopeClaim   string
	ScopeChecker *ScopeChecker
}

func NewJWKUserService(config JWKAuthServiceConfig) *JWKUserService {
	jwks, err := keyfunc.Get(
		config.KeyURL,
		keyfunc.Options{
			RefreshInterval: config.RefreshTime,
			RefreshErrorHandler: func(err error) {
				log.WithFields(
					log.Fields{
						"service": "JWKUserService",
						"error":   err,
					},
				).Error("Failed to refresh JWKs")
			},
		},
	)
	if err != nil {
		log.WithFields(
			log.Fields{
				"service": "JWKUserService",
				"error":   err,
			},
		).Panic("Failed to init JWKs")
	}

	return &JWKUserService{
		jwks: jwks,
		roleGetter: &RoleGetter{
			Default: string(entity.UserRoleUnknown),
			OnRoles: []OnRole{
				{
					OnRole:     config.SuperAdminRole,
					ReturnRole: entity.UserRoleSuperAdmin.String(),
				},
				{
					OnRole:     config.AdminRole,
					ReturnRole: entity.UserRoleAdmin.String(),
				},
				{
					OnRole:     config.UserRole,
					ReturnRole: entity.UserRoleUser.String(),
				},
			},
		},
		ScopeChecker: &ScopeChecker{
			RequiredScope: config.RequiredScope,
		},
		RoleClaim:  config.RoleClaim,
		ScopeClaim: config.ScopeClaim,
	}
}

/*
AuthUser checks user's token and returns user's info

bearerToken - user's token with bearer prefix

throws errors:

 1. wrapped ErrTokenNotValid with reasons:

    1.1 ErrTokenExpired - if token expired

    1.2 ErrFailedToParseToken - if token not valid

    1.3 uknown error

 2. ErrDontHaveScope - if token don't have needed scope

 3. ErrDontHaveRole - if token don't have needed role

 4. ErrDontHaveUserID - if token don't have user id in specifed claim or in client_id
*/
func (s *JWKUserService) AuthUser(
	ctx context.Context,
	// JWT token with bearer prefix
	bearerToken string,
) (aggregate.User, error) {
	token := strings.ReplaceAll(bearerToken, "Bearer ", "")

	var claim jwt.MapClaims

	parsed, err := jwt.ParseWithClaims(token, &claim, s.jwks.Keyfunc)
	if validErr, ok := err.(*jwt.ValidationError); ok {
		switch validErr.Errors {
		case jwt.ValidationErrorExpired:
			return aggregate.User{}, errors.Wrap(ErrTokenExpired, ErrTokenNotValid)
		default:
			return aggregate.User{}, errors.Wrap(validErr, ErrTokenNotValid)
		}
	} else if err != nil {
		return aggregate.User{}, errors.Wrap(ErrFailedToParseToken, ErrTokenNotValid)
	}

	if !parsed.Valid {
		return aggregate.User{}, errors.Wrap(ErrFailedToParseToken, ErrTokenNotValid)
	}

	if err := s.CheckScope(claim[s.ScopeClaim]); err != nil {
		return aggregate.User{}, err
	}

	role, err := s.GetRole(claim[s.RoleClaim])
	if err != nil {
		return aggregate.User{}, err
	}

	// MARK: ignore error, while we not decide how to hand other service auth
	userId, _ := s.GetUserID(claim)

	user, err := aggregate.NewUser(
		fmt.Sprint(userId),
		role,
	)
	if err != nil {
		return aggregate.User{}, err
	}

	return user, nil
}

func (s *JWKUserService) GetUserID(claims jwt.MapClaims) (string, error) {
	// TODO: add client_id claim if needed; wait answer in slack
	if userId, find := claims["sub"]; find {
		return fmt.Sprint(userId), nil
	}

	return "", ErrDontHaveUserID
}

func (s *JWKUserService) GetRole(claims interface{}) (entity.UserRole, error) {
	var strRoles []string
	{
		switch v := claims.(type) {
		case string:
			strRoles = []string{v}
		case []any:
			lo.ForEach(
				v,
				func(role any, _ int) {
					strRoles = append(strRoles, fmt.Sprint(role))
				},
			)
		case any:
			strRoles = []string{fmt.Sprint(v)}
		}
	}

	role := entity.UserRoleFromString(s.roleGetter.GetRole(strRoles))
	if role == entity.UserRoleUnknown {
		return entity.UserRoleUnknown, ErrDontHaveRole
	}

	return role, nil
}

func (s *JWKUserService) CheckScope(claims interface{}) error {
	var strScopes []string
	{
		switch v := claims.(type) {
		case string:
			strScopes = []string{v}
		case []any:
			lo.ForEach(
				v,
				func(scope any, _ int) {
					strScopes = append(strScopes, fmt.Sprint(scope))
				},
			)
		case any:
			strScopes = []string{fmt.Sprint(v)}
		}
	}

	if !s.ScopeChecker.CheckScope(strScopes) {
		return ErrDontHaveScope
	}

	return nil
}
