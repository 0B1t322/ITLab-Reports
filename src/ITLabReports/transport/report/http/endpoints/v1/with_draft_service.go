package endpoints

import (
	"context"

	agragate "github.com/RTUITLab/ITLab-Reports/aggragate/report"
	"github.com/RTUITLab/ITLab-Reports/pkg/endpoint"
	"github.com/RTUITLab/ITLab-Reports/transport/report/http/dto/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/report/http/errors"
)

type DraftService interface {
	IsDraftNotFoundErr(error) bool
	IsDraftIdNotValidErr(error) bool
	// Should throws errors
	// 	Draft not found
	// 	Draft id is invalid
	GetDraft(ctx context.Context, id string) (*agragate.Report, error)

	// Should throws errors
	// 	Draft not found
	// 	Draft id is invalid
	DeleteDraft(ctx context.Context, id string) error
}

type CreateReportFromDraft = endpoint.Endpoint[*dto.CreateReportFromDraftReq, *dto.CreateReportFromDraftResp]

type DraftServiceEndpoints struct {
	CreateReportFromDraft CreateReportFromDraft
}

func NewDraftServiceEndpoints(
	s DraftService,
	e Endpoints,
) DraftServiceEndpoints {
	return DraftServiceEndpoints{
		CreateReportFromDraft: makeDraftServiceEndpoints(s, e),
	}
}


func makeDraftServiceEndpoints(
	s DraftService,
	e Endpoints,
) CreateReportFromDraft {
	return func(
		ctx context.Context,
		req *dto.CreateReportFromDraftReq,
	) (responce *dto.CreateReportResp, err error) {
		r, err := s.GetDraft(
			ctx,
			req.ID,
		)
		if s.IsDraftNotFoundErr(err) {
			return nil, errors.DraftNotFound
		} else if s.IsDraftIdNotValidErr(err) {
			return nil, errors.DraftIdNotValud
		} else if err != nil {
			return nil, err
		}

		resp, err := e.CreateReport(
			ctx,
			&dto.CreateReportReq{
				Name: &r.Report.Name,
				Text: r.Report.Text,
				Implementor: r.Assignees.Implementer,
				Reporter: r.Assignees.Reporter,
			},
		)
		if err != nil {
			return nil, err
		}

		if err := s.DeleteDraft(
			ctx,
			req.ID,
		); err != nil {
			return nil, err
		}

		return resp, nil
	}
}
