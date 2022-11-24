package reports

import (
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/grpc/shared/errors"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/grpc/shared/utils"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/grpc/shared/views"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/shared/reports"
	shared "github.com/RTUITLab/ITLab-Reports/internal/controllers/shared/utils"
)

// Views mappers
var (
	ReportFrom = views.ReportFrom
)

var (
	TokenFromContext = utils.TokenFromContext
	ProtoSortOrderTo = utils.ProtoSortOrderTo
	PaidStateTo      = utils.PaidStateTo
	HasMore          = shared.HasMore
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
	AuthErrorsHandler = errors.AuthErrorsHandler
	ErrorFormatter    = errors.ErrorFormatter
)
