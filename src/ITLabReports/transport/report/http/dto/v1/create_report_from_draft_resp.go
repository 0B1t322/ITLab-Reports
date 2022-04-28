package dto

import (
	"context"
	"encoding/json"
	"net/http"
)

type CreateReportFromDraftResp = CreateReportResp


func EncodeCreateReportFromDraftResp(
	ctx context.Context,
	w http.ResponseWriter,
	resp *CreateReportFromDraftResp,
) error {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resp)
}