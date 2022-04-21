package dto

import (
	"context"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"

	"github.com/gorilla/mux"
)

type GetReportReq struct {
	ID string
}

func (g *GetReportReq) ToEndopointReq() *reqresp.GetReportReq {
	return &reqresp.GetReportReq{
		ID: g.ID,
	}
}

func DecodeGetReportReq(
	ctx		context.Context,
	r		*http.Request,
) (*GetReportReq, error) {
	vars := mux.Vars(r)

	req := &GetReportReq{
		ID: vars["id"],
	}

	return req, nil
}