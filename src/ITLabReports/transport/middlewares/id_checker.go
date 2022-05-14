package middlewares

import (
	"fmt"

	"github.com/RTUITLab/ITLab-Reports/pkg/errors"
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares/context"
)

var (
	// Return if err is incorrect
	ErrIncorectId = errors.New("ID is incorrect")

	// If error is uncatchable
	ErrFaieldToValidateId = errors.New("Failed to validate id")
)

type IdChecker interface {
	// if return error, should error 
	// Where will be information about incorrect id
	CheckIds(token string, ids []string) error

	// Parse error and get incorrect id
	GetIncorrectId(error) string

	// Say that error is cathable
	IsIncorrectIdError(error) bool
}

type ReqWithCheckabeIds interface {
	GetIds() []string
}

func CheckIds[Req ReqWithCheckabeIds, Resp any](
	checker IdChecker,
) MiddlewareWithContext[Req, Resp] {
	return func(next EndpointWithContext[Req, Resp]) EndpointWithContext[Req, Resp] {
		return func(
			ctx context.MiddlewareContext,
			request Req,
		) (Resp, error){
			token, err := ctx.GetToken()
			if err != nil {
				return *new(Resp), TokenNotValid
			}

			ids := request.GetIds()
			if len(ids) > 0 {
				err = checker.CheckIds(token, ids)
				if checker.IsIncorrectIdError(err) {
					return *new(Resp), errors.Wrap(fmt.Errorf("%s", checker.GetIncorrectId(err)), ErrIncorectId)
				} else if err != nil {
					return *new(Resp), errors.Wrap(err, ErrFaieldToValidateId)
				}
			}

			return next(ctx, request)
		}
	}
}