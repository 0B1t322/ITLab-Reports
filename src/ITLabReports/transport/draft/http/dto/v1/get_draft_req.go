package dto

import (
	"context"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
	"github.com/gorilla/mux"
)

type GetDraftReq struct {
	ID string `json:"-"`
}

func (g *GetDraftReq) GetID() string {
	return g.ID
}

func (g *GetDraftReq) ToEndopointReq() *reqresp.GetReportReq {
	return &reqresp.GetReportReq{
		ID: g.ID,
	}
}

func DecodeGetDraftReq(
	ctx context.Context,
	r *http.Request,
) (*GetDraftReq, error) {
	vars := mux.Vars(r)

	return &GetDraftReq{
		ID: vars["id"],
	}, nil
}