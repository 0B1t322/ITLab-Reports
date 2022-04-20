package reqresp

import (
	reportdomain "github.com/RTUITLab/ITLab-Reports/domain/report"
)

type CountReportsReq struct {
	Params *reportdomain.GetReportsFilterFieldsWithOrAnd
}

type CountReportsResp struct {
	Count int64
}