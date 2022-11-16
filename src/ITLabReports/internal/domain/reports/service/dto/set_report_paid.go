package dto

import "github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"

type SetReportPaidReq struct {
	ID string
	By aggregate.User
}
