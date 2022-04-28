package dto

import (
	"context"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/domain/report"
	"github.com/RTUITLab/ITLab-Reports/pkg/optional"
	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
	"github.com/clarketm/json"
	"github.com/gorilla/mux"
)

type UpdateDraftReq struct {
	ID          string                    `json:"-" swaggerignore:"true"`
	Name        optional.Optional[string] `json:"name" swaggertype:"string" extensions:"x-nullable"`
	Text        optional.Optional[string] `json:"text" swaggertype:"string" extensions:"x-nullable"`
	Implementer optional.Optional[string] `json:"implementer" swaggertype:"string" extensions:"x-nullable"`
}

func (u *UpdateDraftReq) GetID() string {
	return u.ID
}

func (u *UpdateDraftReq) ToEndpointReq() *reqresp.UpdateReportReq {
	return &reqresp.UpdateReportReq{
		ID: u.ID,
		Params: report.UpdateReportParams{
			Name:        u.Name,
			Text:        u.Text,
			Implementer: u.Implementer,
		},
	}
}

func DecodeUpdateDraftReq(
	ctx context.Context,
	r *http.Request,
) (*UpdateDraftReq, error) {
	vars := mux.Vars(r)

	req := &UpdateDraftReq{
		ID: vars["id"],
	}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}

	return req, nil
}
