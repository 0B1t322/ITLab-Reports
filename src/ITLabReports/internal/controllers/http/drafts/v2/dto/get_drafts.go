package dto

import "github.com/RTUITLab/ITLab-Reports/internal/controllers/http/shared/view"

type GetDraftsReq struct {
	Limit  int64 `json:"-" form:"limit"  validate:"required" swaggerignore:"true"`
	Offset int64 `json:"-" form:"offset" validate:"required" swaggerignore:"true"`
}

type GetDraftsResp struct {
	Count      int64            `json:"count"`
	Drafts     []view.DraftView `json:"items"`
	Offset     int64            `json:"offset"`
	Limit      int64            `json:"limit"`
	TotalCount int64            `json:"total_count"`
	HasMore    bool             `json:"has_more"`
}
