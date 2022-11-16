package reports

import (
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/http/reports/v2/dto"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/http/shared/errors"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/http/shared/utils"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/http/shared/view"
)

type (
	GetReportsReq  = dto.GetReportsReq
	GetReportsResp = dto.GetReportsResp
)

type (
	ReportView = view.ReportView
)

var (
	ReportViewFrom  = view.ReportViewFrom
	ReportsViewFrom = view.ReportsViewFrom
	HasMore         = utils.HasMore
)

type (
	ErrorFormatter   = errors.ErrorFormatter
	AuthErrorHandler = errors.AuthErrorHandler
)
