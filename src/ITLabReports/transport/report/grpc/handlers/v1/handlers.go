package handlers

import (
	"github.com/RTUITLab/ITLab-Reports/pkg/transport/grpc"
	"github.com/RTUITLab/ITLab-Reports/transport/grpc/options/serverbefore"
	"github.com/RTUITLab/ITLab-Reports/transport/report/grpc/dto/v1"
	. "github.com/RTUITLab/ITLab-Reports/transport/report/grpc/endpoints/v1"
	gt "github.com/go-kit/kit/transport/grpc"
)

func GetReportHandler(
	e Endpoints,
) gt.Handler {
	return grpc.NewServer(
		e.GetReport,
		dto.DecodeGetReportReq,
		dto.EncodeGetReportResp,
		gt.ServerBefore(
			serverbefore.TokenFromReq,
		),
	)
}

func GetReportImplementerHandler(
	e Endpoints,
) gt.Handler {
	return grpc.NewServer(
		e.GetReportImplementer,
		dto.DecodeGetReportImplementerReq,
		dto.EncodeGetReportImplementerResp,
		gt.ServerBefore(
			serverbefore.TokenFromReq,
		),
	)
}

 func GetReportsHandler(
	e Endpoints,
) gt.Handler {
	return grpc.NewServer(
		e.GetReports,
		dto.DecodeGetReportsListReq,
		dto.EncodeGetReportsResp,
		gt.ServerBefore(
			serverbefore.TokenFromReq,
		),
	)
}

 func GetReportsPaginatedHandler(
	e Endpoints,
) gt.Handler {
	return grpc.NewServer(
		e.GetReportsPaginated,
		dto.DecodeGetReportsPaginatedReq,
		dto.EncodeGetReportsPaginatedResp,
		gt.ServerBefore(
			serverbefore.TokenFromReq,
		),
	)
}