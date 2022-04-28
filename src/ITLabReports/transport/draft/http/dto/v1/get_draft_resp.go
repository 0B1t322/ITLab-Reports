package dto

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
)

type GetDraftResp struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Text      string          `json:"text"`
	Date      string          `json:"date"`
	Assignees GetAssignesResp `json:"assignees"`
}

func GetDraftRespFrom(
	resp *reqresp.GetReportResp,
) *GetDraftResp {
	return &GetDraftResp{
		ID:   resp.Report.GetID(),
		Name: resp.Report.GetName(),
		Text: resp.Report.GetText(),
		Date: resp.Report.GetDateString(),
		Assignees: GetAssignesResp{
			Reporter:    resp.GetReporter(),
			Implementer: resp.GetImplementer(),
		},
	}
}

func EncodeGetDraftResp(
	ctx context.Context,
	w http.ResponseWriter,
	resp *GetDraftResp,
) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp)
}
