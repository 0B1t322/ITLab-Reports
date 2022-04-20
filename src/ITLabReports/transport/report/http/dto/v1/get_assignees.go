package dto

import (
	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
)

type GetAssigneesResp struct {
	Reporter string	`json:"reporter"`
	Implementer string `json:"implementer"`
}

func GetAssignesRespFrom(resp *reqresp.GetReportResp) GetAssigneesResp {
	r := resp.Report
	return GetAssigneesResp{
		Reporter: r.GetReporter(),
		Implementer: r.GetImplementer(),
	}
}