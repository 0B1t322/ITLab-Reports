package dto

import "github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"

type GetReportReq struct {
	ID string
	By aggregate.User
}
