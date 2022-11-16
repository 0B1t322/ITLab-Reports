package dto

import "github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"

type CreateReportFromDraftReq struct {
	DraftID string
	By      aggregate.User
}
