package dto

import "github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"

type GetDraftReq struct {
	ID string
	By aggregate.User
}
