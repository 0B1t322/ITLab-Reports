package handlers

import (
	"github.com/RTUITLab/ITLab-Reports/pkg/transport/grpc"
	"github.com/RTUITLab/ITLab-Reports/transport/draft/grpc/dto/v1"
	. "github.com/RTUITLab/ITLab-Reports/transport/draft/grpc/endpoints/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/grpc/options/serverbefore"
	gt "github.com/go-kit/kit/transport/grpc"
)

func GetDraftHandler(
	e Endpoints,
) gt.Handler {
	return grpc.NewServer(
		e.GetDraft,
		dto.DecodeGetDraftReq,
		dto.EncodeGetDraftResp,
		gt.ServerBefore(
			serverbefore.TokenFromReq,
		),
	)
}

func GetDraftsHandler(
	e Endpoints,
) gt.Handler {
	return grpc.NewServer(
		e.GetDrafts,
		dto.DecodeGetDraftsReq,
		dto.EncodeGetDraftsResp,
		gt.ServerBefore(
			serverbefore.TokenFromReq,
		),
	)
}

func GetDraftsPaginatedHandler(
	e Endpoints,
) gt.Handler {
	return grpc.NewServer(
		e.GetDraftsPaginated,
		dto.DecodeGetDraftsPaginatedReq,
		dto.EncodeGetDraftsPaginatedResp,
		gt.ServerBefore(
			serverbefore.TokenFromReq,
		),
	)
}
