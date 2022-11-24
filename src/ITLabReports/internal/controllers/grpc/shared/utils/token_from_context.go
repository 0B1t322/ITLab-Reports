package utils

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func TokenFromContext(ctx context.Context) (token string) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return token
	}

	header := md.Get("Authorization")
	if len(header) == 0 {
		return token
	}

	return header[0]
}
