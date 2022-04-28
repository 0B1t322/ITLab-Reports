package endpoints

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/pkg/endpoint"
	"github.com/RTUITLab/ITLab-Reports/service/reports"
	"github.com/RTUITLab/ITLab-Reports/transport/draft/http/dto/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
	"github.com/RTUITLab/ITLab-Reports/transport/draft/http/errors"
)

type GetDraft = endpoint.Endpoint[*dto.GetDraftReq, *dto.GetDraftResp]
type CreateDraft = endpoint.Endpoint[*dto.CreateDraftReq, *dto.CreateDraftResp]
type DeleteDraft = endpoint.Endpoint[*dto.DeleteDraftReq, *dto.DeleteDraftResp]
type UpdateDraft = endpoint.Endpoint[*dto.UpdateDraftReq, *dto.UpdateDraftResp]
type GetDrafts = endpoint.Endpoint[*dto.GetDraftsReq, *dto.GetDraftsResp]

type Endpoints struct {
	GetDraft    GetDraft
	CreateDraft CreateDraft
	DeleteDraft DeleteDraft
	UpdateDraft UpdateDraft
	GetDrafts   GetDrafts
}

func NewEndpoints(
	e report.Endpoints,
) Endpoints {
	return Endpoints{
		GetDraft: makeGetDraft(e),
		CreateDraft: makeCreateDraft(e),
		DeleteDraft: makeDeleteDraft(e),
		UpdateDraft: makeUpdateDraft(e),
		GetDrafts: makeGetDrafts(e),
	}
}

func makeGetDraft(
	e report.Endpoints,
) GetDraft {
	return func(
		ctx context.Context, 
		req *dto.GetDraftReq,
	) (*dto.GetDraftResp, error) {
		resp, err := e.GetReport(
			ctx,
			req.ToEndopointReq(),
		)
		switch {
		case err == reports.ErrReportIDNotValid:
			return nil, errors.DraftIDIsInvalid
		case err == reports.ErrReportNotFound:
			return nil, errors.DraftNotFound
		case err != nil:
			return nil, err
		}

		return dto.GetDraftRespFrom(resp), nil
	}
}

func makeCreateDraft(
	e report.Endpoints,
) CreateDraft {
	return func(
		ctx context.Context, 
		req *dto.CreateDraftReq,
	) (responce *dto.CreateDraftResp, err error) {
		resp, err := e.CreateReport(
			ctx,
			req.ToEndpointReq(),
		)
		if err != nil {
			return nil, err
		}

		return dto.CreateDraftRespFrom(resp), nil
	}
}

func makeDeleteDraft(
	e report.Endpoints,
) DeleteDraft {
	return func(
		ctx context.Context, 
		req *dto.DeleteDraftReq,
	) (responce *dto.DeleteDraftResp, err error) {
		_, err = e.DeleteReport(
			ctx,
			req.ToEndopointReq(),
		)
		switch {
		case err == reports.ErrReportIDNotValid:
			return nil, errors.DraftIDIsInvalid
		case err == reports.ErrReportNotFound:
			return nil, errors.DraftNotFound
		case err != nil:
			return nil, err
		}

		return &dto.DeleteDraftResp{}, nil
	}
}

func makeUpdateDraft(
	e report.Endpoints,
) UpdateDraft {
	return func(
		ctx context.Context, 
		req *dto.UpdateDraftReq,
	) (responce *dto.GetDraftResp, err error) {
		resp, err := e.UpdateReport(
			ctx,
			req.ToEndpointReq(),
		)
		switch {
		case err == reports.ErrReportIDNotValid:
			return nil, errors.DraftIDIsInvalid
		case err == reports.ErrReportNotFound:
			return nil, errors.DraftNotFound
		case err != nil:
			return nil, err
		}


		return dto.UpdateDraftRespFrom(resp), nil
	}
}

func makeGetDrafts(
	e report.Endpoints,
) GetDrafts {
	return func(
		ctx context.Context, 
		req *dto.GetDraftsReq,
	) (responce *dto.GetDraftsResp, err error) {
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