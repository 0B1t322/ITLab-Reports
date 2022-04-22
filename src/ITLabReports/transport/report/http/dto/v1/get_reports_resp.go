package dto

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
)

type GetReportsResp []*GetReportResp

func GetReportsRespFrom(resp *reqresp.GetReportsResp) GetReportsResp {
	reports := resp.Reports
	httpResp := []*GetReportResp{}

	for _, r := range reports {
		httpResp = append(httpResp, GetReportRespFrom(&reqresp.GetReportResp{Report: r}))
	}

	return httpResp
}

func EncodeGetReportsResp(
	ctx context.Context,
	w http.ResponseWriter,
	resp *GetReportsResp,
) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resp)
}