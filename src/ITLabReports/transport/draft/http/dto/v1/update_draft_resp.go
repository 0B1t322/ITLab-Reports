package dto

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
)

type UpdateDraftResp = GetDraftResp

func UpdateDraftRespFrom(
	resp *reqresp.UpdateReportResp,
) *UpdateDraftResp {
	return GetDraftRespFrom(&reqresp.GetReportResp{Report: resp.Report})
}

func EncodeUpdateDraftResp(
	ctx context.Context,
	w http.ResponseWriter,
	resp *UpdateDraftResp,
) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp)
}