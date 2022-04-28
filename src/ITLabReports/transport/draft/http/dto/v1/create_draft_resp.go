package dto

import (
	"context"
	"net/http"

	v1 "github.com/RTUITLab/ITLab-Reports/transport/report/http/dto/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
	"github.com/clarketm/json"
)

type GetAssignesResp = v1.GetAssigneesResp

type CreateDraftResp struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Text     string          `json:"text"`
	Date     string          `json:"date"`
	Assignes GetAssignesResp `json:"assignees"`
}

func CreateDraftRespFrom(
	resp *reqresp.CreateReportResp,
) *CreateDraftResp {
	return &CreateDraftResp{
		ID:       resp.Report.GetID(),
		Name:     resp.Report.GetName(),
		Text:     resp.Report.GetText(),
		Date:     resp.Report.GetDateString(),
		Assignes: GetAssignesResp{
			Reporter: resp.Report.Assignees.Reporter,
			Implementer: resp.Report.Assignees.Implementer,
		},
	}
}

func EncodeCreateDraftResp(
	ctx context.Context,
	w http.ResponseWriter,
	resp *CreateDraftResp,
) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(resp)
}