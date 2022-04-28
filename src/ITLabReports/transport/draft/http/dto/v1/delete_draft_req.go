package dto

import (
	"context"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
	"github.com/gorilla/mux"
)

type DeleteDraftReq struct {
	ID string	`json:"-"`
}

func (d *DeleteDraftReq) GetID() string {
	return d.ID
}

func (d *DeleteDraftReq) ToEndopointReq() *reqresp.DeleteReportReq {
	return &reqresp.DeleteReportReq{
		ID: d.ID,
	}
}

func DecodeDeleteDraftReq(
	ctx context.Context,
	r *http.Request,
) (*DeleteDraftReq, error) {
	vars := mux.Vars(r)

	return &DeleteDraftReq{
		ID: vars["id"],
	}, nil
}