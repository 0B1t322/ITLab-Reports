package utils

import (
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	pb "github.com/RTUITLab/ITLab/proto/reports/v1"
)

func PaidStateTo(paidState pb.GetReportsPaginatedReq_FilterParams_PaidState) aggregate.ReportState {
	switch paidState {
	case pb.GetReportsPaginatedReq_FilterParams_NOT_PAID:
		return aggregate.ReportStateCreated
	default:
		return aggregate.ReportStatePaid
	}
}
