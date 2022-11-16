package reports

import (
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/http/reports/v1/dto"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/http/shared/errors"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/http/shared/view"
)

// Alias dto
type (
	CreateReportReq          = dto.CreateReportReq
	GetReportReq             = dto.GetReportReq
	GetReportsReq            = dto.GetReportsReq
	GetEmployeeReportsReq    = dto.GetEmployeeReportsReq
	CreateReportFromDraftReq = dto.CreateReportFromDraftReq
)

// Alias View
type (
	ReportView = view.ReportView
)

// Alias View marshalling
var (
	ReportViewFrom  = view.ReportViewFrom
	ReportsViewFrom = view.ReportsViewFrom
)

// Alias errors handlers
type (
	ErrorFormatter   = errors.ErrorFormatter
	AuthErrorHandler = errors.AuthErrorHandler
)
