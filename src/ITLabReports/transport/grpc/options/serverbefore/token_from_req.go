package serverbefore

import (
	"context"

	"google.golang.org/grpc/metadata"
	mcontext "github.com/RTUITLab/ITLab-Reports/transport/middlewares/context"
)

func TokenFromReq(
	ctx context.Context, 
	m metadata.MD,
) context.Context {
	mctx := mcontext.New(ctx)
	tokenHeader := m.Get("Authorization")
	if len(tokenHeader) == 0 {
		return mctx
	}

	mctx.SetToken(tokenHeader[0])
	return mctx
}