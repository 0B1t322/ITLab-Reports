package salary

import (
	"context"
	"fmt"

	"github.com/RTUITLab/ITLab-Reports/pkg/optional"
	services "github.com/RTUITLab/ITLab/proto/salary/v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ExternalGRPCSalaryService struct {
	cleint services.ApprovedReportsSalaryClient
}

func NewExternalGRPCSalaryService(
	cleint services.ApprovedReportsSalaryClient,
) SalaryService {
	return &ExternalGRPCSalaryService{
		cleint: cleint,
	}
}

func (e *ExternalGRPCSalaryService) GetApprovedReportsIds(
	ctx context.Context,
	token string,
	userId optional.Optional[string],
) ([]string, error) {
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("Authorization", token))

	resp, err := e.cleint.GetApprovedReportsId(
		ctx,
		&services.GetApprovedReportsIdReq{
			UserId: userId.Value,
		},
	)
	if err != nil {
		s, _ := status.FromError(err)
		return nil, fmt.Errorf("Unexpected responce status code: %v, message: %v", s.Code(), s.Message())
	}

	return resp.ReportId, nil
}