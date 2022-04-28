package dto

import (
	"context"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
	"github.com/clarketm/json"
)

type GetDraftsResp struct {
	Drafts []*GetDraftResp `json:"drafts"`
}

func GetDraftsRespFrom(from *reqresp.GetReportsResp) *GetDraftsResp {
	resp := &GetDraftsResp{}

	for _, report := range from.Reports {
		resp.Drafts = append(resp.Drafts, GetDraftRespFrom(&reqresp.GetReportResp{Report: report}))
	}

	return resp
}

func EncodeGetDraftsResp(
	ctx context.Context,
	w http.ResponseWriter,
	resp *GetDraftsResp,
) error {
	return json.NewEncoder(w).Encode(resp)
}