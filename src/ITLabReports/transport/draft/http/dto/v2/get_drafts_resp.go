package dto

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/transport/draft/http/dto/v1"
)

type GetDraftsResp struct {
	Count       int                 `json:"count"`
	HasMore     bool                `json:"hasMore"`
	Items       []*dto.GetDraftResp `json:"items"`
	Limit       int                 `json:"limit"`
	Offset      int                 `json:"offset"`
	TotalResult int                 `json:"totalResult"`
} //@name dto.GetDraftRespV2

func EncodeGetDraftResp(
	ctx context.Context,
	w http.ResponseWriter,
	resp *GetDraftsResp,
) error {
	if resp.Items == nil {
		resp.Items = []*dto.GetDraftResp{}
	}

	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resp)
}
