package dto

import drafts "github.com/RTUITLab/ITLab-Reports/internal/domain/drafts/repository"

type GetDraftsReq struct {
	Query drafts.GetDraftsQuery
}
