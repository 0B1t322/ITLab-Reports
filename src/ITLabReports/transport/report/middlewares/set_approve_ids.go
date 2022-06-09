package middlewares

import (
	"github.com/RTUITLab/ITLab-Reports/pkg/errors"
	. "github.com/RTUITLab/ITLab-Reports/transport/middlewares"
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares/context"
)

var (
	ErrFailedToGetApprovedReportsIds = errors.New("Failed to get approved reports ids")
)

type ReqWithSetApprovedReportsIds interface {
	IsOnlyApprovedReports() bool
	SetOnlyApprovedReports(ids ...string)
}

type ApprovedReportsIdsGetter interface {
	GetApprovedReportsIdsForUser(userId string, token string) ([]string, error)
	GetApprovedReportsIds(token string) ([]string, error)
}

type RoleGetter interface {
	GetAdminRole() string

	GetUserRole() string

	GetSuperAdminRole() string
}

func SetApprovedReportsIds[Req ReqWithSetApprovedReportsIds, Resp any](
	idsGetter ApprovedReportsIdsGetter,
	roleGetter RoleGetter,
) MiddlewareWithContext[Req, Resp] {
	return func(next EndpointWithContext[Req, Resp]) EndpointWithContext[Req, Resp] {
		return func(
			ctx context.MiddlewareContext,
			request Req,
		) (Resp, error) {
			if !request.IsOnlyApprovedReports() {
				return next(ctx, request)
			}

			token, err := ctx.GetToken()
			if err != nil {
				return next(ctx, request)
			}

			role, err := ctx.GetRole()
			if err != nil {
				return next(ctx, request)
			}

			var ids []string
			{
				if role == roleGetter.GetUserRole() {
					userId, _ := ctx.GetUserID()
					ids, err = idsGetter.GetApprovedReportsIdsForUser(
						userId,
						token,
					)
				} else if role == roleGetter.GetAdminRole() || role == roleGetter.GetSuperAdminRole() {
					ids, err = idsGetter.GetApprovedReportsIds(
						token,
					)
				}
			}
			if err != nil {
				return *new(Resp), errors.Wrap(err, ErrFailedToGetApprovedReportsIds)
			}
			request.SetOnlyApprovedReports(ids...)

			return next(ctx, request)
		}
	}
}