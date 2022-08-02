package dto

import (
	"context"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/domain/report"
	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
	"github.com/clarketm/json"
	"github.com/gorilla/mux"
	"github.com/samber/mo"
)

type UpdateDraftReq struct {
	ID          string                    `json:"-" swaggerignore:"true"`
	Name        string `json:"name" swaggertype:"string" extensions:"x-nullable"`
	Text        string `json:"text" swaggertype:"string" extensions:"x-nullable"`
	Implementer string `json:"implementer" swaggertype:"string" extensions:"x-nullable"`
}

func (u *UpdateDraftReq) GetIds() []string {
	var ids []string
	{
		if u.Implementer != "" {
			ids = append(ids, u.Implementer)
		}
	}
	return ids
}

func (u *UpdateDraftReq) GetID() string {
	return u.ID
}

func (u *UpdateDraftReq) ToEndpointReq() *reqresp.UpdateReportReq {
	name := mo.None[string]()
	if u.Name != "" {
		name = mo.Some(u.Name)
	}

	text := mo.None[string]()
	if u.Text != "" {
		text = mo.Some(u.Text)
	}

	implementer := mo.None[string]()
	if u.Implementer != "" {
		implementer = mo.Some(u.Implementer)
	}

	return &reqresp.UpdateReportReq{
		ID: u.ID,
		Params: report.UpdateReportParams{
			Name:        name,
			Text:        text,
			Implementer: implementer,
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
