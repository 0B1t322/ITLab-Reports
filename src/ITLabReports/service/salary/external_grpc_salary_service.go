package salary

import (
	"context"
	"fmt"

	services "github.com/RTUITLab/ITLab/proto/salary/v1"
	"github.com/samber/mo"
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
	UserId mo.Option[string],
) ([]string, error) {
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("Authorization", token))
	var userId *string = nil
	{
		if UserId.IsPresent() {
			value := UserId.MustGet()
			userId = &value
		}
	}
	resp, err := e.cleint.GetApprovedReportsId(
		ctx,
		&services.GetApprovedReportsIdReq{
			UserId: userId,
		},
	)
	if err != nil {
		s, _ := status.FromError(err)
		return nil, fmt.Errorf("Unexpected responce status code: %v, message: %v", s.Code(), s.Message())
	}

	return resp.ReportId, nil
}