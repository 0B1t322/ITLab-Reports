package dto

import (
	"context"
	"encoding/json"
	"net/http"

	v1 "github.com/RTUITLab/ITLab-Reports/transport/report/http/dto/v1"
)

type GetReportsResp struct {
	Count       int              `json:"count"`
	HasMore     bool             `json:"hasMore"`
	Items       []*v1.GetReportResp `json:"items"`
	Limit       int              `json:"limit"`
	Offset      int              `json:"offset"`
	TotalResult int              `json:"totalResult"`
}//@name dto.GetReportsRespV2

func EncodeGetReportsResp(
	ctx context.Context,
	w http.ResponseWriter,
	resp *GetReportsResp,
) error {
	if resp.Items == nil {
		resp.Items = []*v1.GetReportResp{}
	}
	
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resp)
}