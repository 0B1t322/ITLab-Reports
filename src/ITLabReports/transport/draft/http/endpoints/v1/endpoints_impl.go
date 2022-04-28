// Impl of some interfaces
package endpoints

import (
	"context"
	"time"

	agragate "github.com/RTUITLab/ITLab-Reports/aggragate/report"
	"github.com/RTUITLab/ITLab-Reports/entity/assignees"
	"github.com/RTUITLab/ITLab-Reports/entity/report"
	"github.com/RTUITLab/ITLab-Reports/transport/draft/http/dto/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/draft/http/errors"
)

// Implement OwnerChecker
func (e Endpoints) CheckUserIsOwner(
	ctx context.Context,
	userId string,
	draftId string,
) (bool, error) {
	resp, err := e.GetDraft(
		ctx,
		&dto.GetDraftReq{
			ID: draftId,
		},
	)
	if err != nil {
		return false, err
	}

	return resp.Assignees.Reporter == userId, nil
}

type ToDraftServiceAdapter struct {
	Endpoints
}

func ToDraftService(
	e Endpoints,
) ToDraftServiceAdapter {
	return ToDraftServiceAdapter{e}
}

func (t ToDraftServiceAdapter) DeleteDraft(
	ctx context.Context,
	id string,
) error {
	_, err := t.Endpoints.DeleteDraft(
		ctx,
		&dto.DeleteDraftReq{
			ID: id,
		},
	)

	return err
}

func (t ToDraftServiceAdapter) GetDraft(
	ctx context.Context,
	id string,
) (*agragate.Report, error) {
	resp, err := t.Endpoints.GetDraft(
		ctx,
		&dto.GetDraftReq{
			ID: id,
		},
	)
	if err != nil {
		return nil, err
	}

	return &agragate.Report{
		Report: &report.Report{
			ID:   resp.ID,
			Name: resp.Name,
			Date: time.Now().UTC(),
			Text: resp.Text,
		},
		Assignees: &assignees.Assignees{
			Reporter:    resp.Assignees.Reporter,
			Implementer: resp.Assignees.Implementer,
		},
	}, nil
}

func (t ToDraftServiceAdapter) IsDraftNotFoundErr(err error) bool {
	return err == errors.DraftNotFound
}

func (t ToDraftServiceAdapter) IsDraftIdNotValidErr(err error) bool {
	return err == errors.DraftIDIsInvalid
}
