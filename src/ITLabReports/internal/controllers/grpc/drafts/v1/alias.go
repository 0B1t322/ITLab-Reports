package drafts

import (
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/grpc/shared/errors"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/grpc/shared/utils"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/grpc/shared/views"
	shared "github.com/RTUITLab/ITLab-Reports/internal/controllers/shared/utils"
)

// Views mappers
var (
	DraftFrom = views.DraftFrom
)

var (
	TokenFromContext = utils.TokenFromContext
	ProtoSortOrderTo = utils.ProtoSortOrderTo
	PaidStateTo      = utils.PaidStateTo
	HasMore          = shared.HasMore
)

type (
	AuthErrorsHandler = errors.AuthErrorsHandler
	ErrorFormatter    = errors.ErrorFormatter
)
