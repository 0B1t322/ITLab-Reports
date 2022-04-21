package dto

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
)

type GetReportResp struct {
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	Text     string           `json:"text"`
	Date     string           `json:"date"`
	Assignes GetAssigneesResp `json:"assignees"`
}

func (g *GetReportResp) GetImplementer() string {
	return g.Assignes.Implementer
}

func (g *GetReportResp) GetReporter() string {
	return g.Assignes.Reporter
}

func GetReportRespFrom(resp *reqresp.GetReportResp) *GetReportResp {
	r := resp.Report

	return &GetReportResp{
		ID:       r.GetID(),
		Name:     r.GetName(),
		Text:     r.GetText(),
		Date:     r.GetDateString(),
		Assignes: GetAssignesRespFrom(resp),
	}
}

func EncodeGetReportResp(
	ctx context.Context,
	w http.ResponseWriter,
	resp *GetReportResp,
) error {
	return json.NewEncoder(w).Encode(resp)
}
