package endpoints

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/pkg/endpoint"
	"github.com/RTUITLab/ITLab-Reports/service/reports"
	"github.com/RTUITLab/ITLab-Reports/transport/draft/grpc/dto/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/draft/grpc/utils"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
	pb "github.com/RTUITLab/ITLab/proto/reports/v1"
	"github.com/RTUITLab/ITLab/proto/shared"
)

type (
	GetDraft           = endpoint.Endpoint[*dto.GetDraftReq, *dto.GetDraftResp]
	GetDrafts          = endpoint.Endpoint[*dto.GetDraftsReq, *dto.GetDraftsResp]
	GetDraftsPaginated = endpoint.Endpoint[*dto.GetDraftsPaginatedReq, *dto.GetDraftsPaginatedResp]
)

type Endpoints struct {
	GetDraft           GetDraft
	GetDrafts          GetDrafts
	GetDraftsPaginated GetDraftsPaginated
}

func NewEndpoint(
	e report.Endpoints,
) Endpoints {
	return Endpoints{
		GetDraft:           makeGetDraftEndpoint(e),
		GetDrafts:          makeGetDrafts(e),
		GetDraftsPaginated: makeGetDraftsPaginated(e),
	}
}

func makeGetDraftEndpoint(
	e report.Endpoints,
) GetDraft {
	return func(
		ctx context.Context,
		req *dto.GetDraftReq,
	) (response *dto.GetDraftResp, err error) {
		resp, err := e.GetReport(
			ctx,
			req.ToEndpointReq(),
		)
		if err == reports.ErrReportNotFound || err == reports.ErrReportIDNotValid {
			return &dto.GetDraftResp{
				Result: &pb.GetDraftResp_Error{
					Error: pb.DraftsServiceErrors_DRAFT_NOT_FOUND,
				},
			}, nil
		} else if err != nil {
			return nil, err
		}
		return &dto.GetDraftResp{
			Result: &pb.GetDraftResp_Draft{
				Draft: utils.DraftToPBType(resp.Report),
			},
		}, nil
	}
}

func makeGetDrafts(
	e report.Endpoints,
) GetDrafts {
	return func(
		ctx context.Context,
		req *dto.GetDraftsReq,
	) (response *dto.GetDraftsResp, err error) {
		resp, err := e.GetReports(
			ctx,
			req.ToEndpointReq(),
		)
		if err != nil {
			return nil, err
		}

		return dto.GetDraftsRespFrom(resp), nil
	}
}

func makeGetDraftsPaginated(
	e report.Endpoints,
) GetDraftsPaginated {
	return func(
		ctx context.Context,
		req *dto.GetDraftsPaginatedReq,
	) (*dto.GetDraftsPaginatedResp, error) {
		resp, err := e.GetReports(
			ctx,
			req.ToEndpointReq(),
		)
		if err != nil {
			return nil, err
		}

		countDrafts, err := e.CountReports(
			ctx,
			&reqresp.CountReportsReq{
				Params: &req.Params.Filter.GetReportsFilterFieldsWithOrAnd,
			},
		)

		response := &dto.GetDraftsPaginatedResp{}
		{
			paginationInfo := &shared.PaginationInfo{}
			{
				paginationInfo.Offset = 0
				paginationInfo.Limit = 0
				paginationInfo.Count = int64(len(resp.Reports))
				paginationInfo.TotalResult = countDrafts.Count
				if req.Params.Limit.IsPresent() {
					paginationInfo.Limit = req.Params.Limit.MustGet()
				}

				if req.Params.Offset.IsPresent() {
					paginationInfo.Offset = req.Params.Offset.MustGet()
				}

				paginationInfo.HasMore = false
				if paginationInfo.Limit != 0 && paginationInfo.TotalResult-paginationInfo.Offset-paginationInfo.Limit > 0 {
					paginationInfo.HasMore = true
				}
			}
			response.PaginationInfo = paginationInfo

			for _, d := range resp.Reports {
				response.Drafts = append(response.Drafts, utils.DraftToPBType(d))
			}
		}

		return response, nil
	}
}
