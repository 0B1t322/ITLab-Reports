package endpoints

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/pkg/endpoint"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
	v1 "github.com/RTUITLab/ITLab-Reports/transport/report/http/dto/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/report/http/dto/v2"
	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
)

type Endpoints struct {
	GetReports		endpoint.Endpoint[*dto.GetReportsReq, *dto.GetReportsResp]
}

func NewEndpoints(
	e report.Endpoints,
) Endpoints {
	return Endpoints{
		GetReports: makeGetReportsEndpoint(e),
	}
}

func makeGetReportsEndpoint(
	e report.Endpoints,
) endpoint.Endpoint[*dto.GetReportsReq, *dto.GetReportsResp] {
	return func(
		ctx context.Context,
		req *dto.GetReportsReq,
	) (*dto.GetReportsResp, error) {
		reports, err := e.GetReports(
			ctx,
			req.ToEndpointReq(),
		)
		if err != nil {
			return nil, err
		}

		countReport, err := e.CountReports(
			ctx,
			&reqresp.CountReportsReq{
				Params: &req.Query.Params.Filter.GetReportsFilterFieldsWithOrAnd,
			},
		)
		if err != nil {
			return nil, err
		}

		offset := 0
		limit := 0
		count := len(reports.Reports)
		totalResult := countReport.Count
		HasMore := totalResult - int64(count) > 0

		if req.Query.Params.Limit.HasValue() {
			limit = int(req.Query.Params.Limit.MustGetValue())
		}

		if req.Query.Params.Offset.HasValue() {
			offset = int(req.Query.Params.Offset.MustGetValue())
		}

		resp := &dto.GetReportsResp{
			Items: v1.GetReportsRespFrom(reports),
			Offset: offset,
			Limit: limit,
			TotalResult: int(totalResult),
			Count: count,
			HasMore: HasMore,
		}

		return resp, nil

	}
}