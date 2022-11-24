package drafts

import (
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/http/drafts/v2/dto"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/http/shared/errors"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/http/shared/view"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/shared/utils"
)

type (
	GetDraftsReq  = dto.GetDraftsReq
	GetDraftsResp = dto.GetDraftsResp
)

type (
	DraftView = view.DraftView
)

var (
	DraftViewFrom = view.DraftViewFrom
	HasMore       = utils.HasMore
)

type (
	ErrorFormatter   = errors.ErrorFormatter
	AuthErrorHandler = errors.AuthErrorHandler
)
