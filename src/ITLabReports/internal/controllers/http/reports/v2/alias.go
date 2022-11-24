package reports

import (
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/http/reports/v2/dto"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/http/shared/errors"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/http/shared/view"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/shared/reports"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/shared/utils"
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

// approved strategy
type (
	ApprovedStateStrategy = reports.ApprovedStateStrategy
	ApprovedStateReq      = reports.ApprovedStateReq
)

var (
	NewInternalApprovedStrategy = reports.NewInternalApprovedStrategy
)

type (
	ErrorFormatter   = errors.ErrorFormatter
	AuthErrorHandler = errors.AuthErrorHandler
)
