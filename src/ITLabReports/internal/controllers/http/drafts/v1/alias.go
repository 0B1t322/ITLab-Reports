package drafts

import (
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/http/drafts/v1/dto"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/http/shared/errors"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/http/shared/view"
)

type (
	GetDraftReq    = dto.GetDraftReq
	UpdateDraftReq = dto.UpdateDraftReq
	CreateDraftReq = dto.CreateDraftReq
	DeleteDraftReq = dto.DeleteDraftReq
	GetDraftsReq   = dto.GetDraftsReq
)

type (
	DraftView  = view.DraftView
	DraftsView = view.DraftsView
)

var (
	DraftViewFrom  = view.DraftViewFrom
	DraftsViewFrom = view.DraftsViewFrom
)

type (
	ErrorFormatter   = errors.ErrorFormatter
	AuthErrorHandler = errors.AuthErrorHandler
)
