package dto

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type CreateReportFromDraftReq struct {
	// id of draft
	ID string `json:"-"`
}

func DecodeCreateReportFromDraftReq(
	ctx context.Context,
	r *http.Request,
) (*CreateReportFromDraftReq, error) {
	vars := mux.Vars(r)

	return &CreateReportFromDraftReq{
		ID: vars["id"],
	}, nil
}