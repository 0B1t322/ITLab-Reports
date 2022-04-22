package dto

import (
	"context"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
	"github.com/clarketm/json"
)

type CreateReportResp GetReportResp

func CreateReportRespFrom(resp *reqresp.CreateReportResp) *CreateReportResp {
	dtoResp := CreateReportResp(*GetReportRespFrom(&reqresp.GetReportResp{Report: resp.Report}))

	return &dtoResp
}

func EncodeCreateReportResp(
	ctx context.Context,
	w http.ResponseWriter,
	resp *CreateReportResp,
) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resp)
}