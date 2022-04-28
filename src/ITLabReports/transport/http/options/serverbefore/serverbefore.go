package serverbefore

import (
	"context"
	"net/http"

	mcontext "github.com/RTUITLab/ITLab-Reports/transport/middlewares/context"
)

func TokenFromReq(
	ctx context.Context,
	r *http.Request,
) context.Context {
	mctx := mcontext.CreateOrGetFrom(ctx)
	token := r.Header.Get("Authorization")
	if token == "" {
		return mctx
	}

	mctx.SetToken(token)
	return mctx
}