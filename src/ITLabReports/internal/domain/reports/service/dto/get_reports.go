package dto

import (
	reports "github.com/RTUITLab/ITLab-Reports/internal/domain/reports/repository"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
)

type GetReportsReq struct {
	Query reports.GetReportsQuery
	By    aggregate.User
}
