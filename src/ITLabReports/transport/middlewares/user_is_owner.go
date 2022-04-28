package middlewares

import (
	"context"
	"errors"

	mcontext "github.com/RTUITLab/ITLab-Reports/transport/middlewares/context"
)

var (
	YouAreNotOwner = errors.New("You are not owner")
)

type OwnerChecker interface {
	CheckUserIsOwner(
		ctx context.Context,
		userId string,
		entid string,
	) (bool, error)
}

type ReqWithID interface {
	GetID() string
}

func UserIsOwner[Req ReqWithID, Resp any](ownerChecker OwnerChecker) MiddlewareWithContext[Req, Resp] {
	return func(next EndpointWithContext[Req, Resp]) EndpointWithContext[Req, Resp] {
		return func(
			ctx mcontext.MiddlewareContext, 
			request Req,
		) (Resp, error) {
			userId, err := ctx.GetUserID()
			if err != nil {
				return *new(Resp), err
			}

			isOwner, err := ownerChecker.CheckUserIsOwner(
				ctx,
				userId,
				request.GetID(),
			)
			if err != nil {
				return *new(Resp), err
			}

			if !isOwner {
				return *new(Resp), YouAreNotOwner
			}

			return next(ctx, request)
		}
	}
}