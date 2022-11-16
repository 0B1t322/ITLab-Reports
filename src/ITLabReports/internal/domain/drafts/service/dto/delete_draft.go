package dto

import "github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"

type DeleteDraftReq struct {
	ID string
	By aggregate.User
}
