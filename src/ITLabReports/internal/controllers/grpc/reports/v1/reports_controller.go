package reports

import (
	"context"
	"fmt"

	"github.com/RTUITLab/ITLab-Reports/internal/common/errors"
	reports "github.com/RTUITLab/ITLab-Reports/internal/domain/reports/service"
	user "github.com/RTUITLab/ITLab-Reports/internal/domain/user/service"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/RTUITLab/ITLab/proto/reports/types"
	reportsgrpc "github.com/RTUITLab/ITLab/proto/reports/v1"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReportsController struct {
	reportService    reports.ReportsService
	auth             user.UserService
	approvedStrategy ApprovedStateStrategy
	marshaller       ReportsMarshaller

	reportsgrpc.UnimplementedReportsServer

	AuthErrorsHandler
	ErrorFormatter
}

func NewReportsController(
	reportService reports.ReportsService,
	auth user.UserService,
	approvedStrategy ApprovedStateStrategy,
) *ReportsController {
	return &ReportsController{
		reportService:    reportService,
		auth:             auth,
		approvedStrategy: approvedStrategy,
	}
}

func NewReportsControllerFrom(i *do.Injector) (*ReportsController, error) {
	reportService := do.MustInvoke[reports.ReportsService](i)
	auth := do.MustInvoke[user.UserService](i)

	var approvedStrategy ApprovedStateStrategy = NewInternalApprovedStrategy()
	{
		ovverideApprovedStrategy, err := do.Invoke[ApprovedStateStrategy](i)
		if err == nil {
			approvedStrategy = ovverideApprovedStrategy
		}
	}

	return NewReportsController(reportService, auth, approvedStrategy), nil
}

func (rc *ReportsController) Build(server *grpc.Server) {
	reportsgrpc.RegisterReportsServer(server, rc)
}

func (rc *ReportsController) HandlerError(err error) error {
	if err := rc.AuthErrorsHandler.Handle(err); err != nil {
		return err
	}

	logrus.WithFields(
		logrus.Fields{
			"controller": "reports",
			"transport":  "grpc",
		},
	).Error(err)
	return rc.FormatError(codes.Internal, err)
}

// Return report implementer by id
// If report not found return REPORT_NOT_FOUND error
func (rc *ReportsController) GetReportImplementer(
	ctx context.Context,
	req *reportsgrpc.GetReportImplementerReq,
) (*reportsgrpc.GetReportImplementerResp, error) {
	user, err := rc.auth.AuthUser(
		ctx,
		TokenFromContext(ctx),
	)
	if err != nil {
		return nil, rc.HandlerError(err)
	}

	report, err := rc.reportService.GetReport(
		ctx,
		reports.GetReportReq{
			ID: req.ReportId,
			By: user,
		},
	)
	if err == reports.ErrReportNotFound {
		return &reportsgrpc.GetReportImplementerResp{
			Result: &reportsgrpc.GetReportImplementerResp_Error{
				Error: reportsgrpc.ReportsServiceErrors_REPORT_NOT_FOUND,
			},
		}, nil
	} else if errors.Is(err, reports.ErrCantGetReport) {
		return nil, rc.FormatError(codes.PermissionDenied, err)
	} else if err != nil {
		return nil, rc.HandlerError(err)
	}

	return &reportsgrpc.GetReportImplementerResp{
		Result: &reportsgrpc.GetReportImplementerResp_Implementer{
			Implementer: report.Assignees.Implementer,
		},
	}, nil
}

// Return report by id
// If report not found return REPORT_NOT_FOUND error
func (rc *ReportsController) GetReport(
	ctx context.Context,
	req *reportsgrpc.GetReportReq,
) (*reportsgrpc.GetReportResp, error) {
	user, err := rc.auth.AuthUser(
		ctx,
		TokenFromContext(ctx),
	)
	if err != nil {
		return nil, rc.HandlerError(err)
	}

	report, err := rc.reportService.GetReport(
		ctx,
		reports.GetReportReq{
			ID: req.ReportId,
			By: user,
		},
	)
	if err == reports.ErrReportNotFound {
		return &reportsgrpc.GetReportResp{
			Result: &reportsgrpc.GetReportResp_Error{
				Error: reportsgrpc.ReportsServiceErrors_REPORT_NOT_FOUND,
			},
		}, nil
	} else if errors.Is(err, reports.ErrCantGetReport) {
		return nil, rc.FormatError(codes.PermissionDenied, err)
	} else if err != nil {
		return nil, rc.HandlerError(err)
	}

	return &reportsgrpc.GetReportResp{
		Result: &reportsgrpc.GetReportResp_Report{
			Report: ReportFrom(report),
		},
	}, nil
}

// Return reports list without pagaination
func (rc *ReportsController) GetReports(
	ctx context.Context,
	req *reportsgrpc.GetReportsReq,
) (*reportsgrpc.GetReportsResp, error) {
	user, err := rc.auth.AuthUser(
		ctx,
		TokenFromContext(ctx),
	)
	if err != nil {
		return nil, rc.HandlerError(err)
	}

	reports, err := rc.reportService.GetReports(
		ctx,
		rc.marshaller.MarshalGetReportsReq(req, user),
	)
	if err != nil {
		return nil, rc.HandlerError(err)
	}

	return &reportsgrpc.GetReportsResp{
		Reports: lo.Map(
			reports,
			func(r aggregate.Report, _ int) *types.Report {
				return ReportFrom(r)
			},
		),
	}, nil
}

// Return reports list with pagaination
func (rc *ReportsController) GetReportsPaginated(
	ctx context.Context,
	req *reportsgrpc.GetReportsPaginatedReq,
) (*reportsgrpc.GetReportsPaginatedResp, error) {
	user, err := rc.auth.AuthUser(
		ctx,
		TokenFromContext(ctx),
	)
	if err != nil {
		return nil, rc.HandlerError(err)
	}

	r := rc.marshaller.MarshallGetReportsPaginatedReq(req, user)

	if filter := req.GetFilterParams(); filter != nil {
		if state := filter.GetPaidState(); state != reportsgrpc.GetReportsPaginatedReq_FilterParams_ALL {
			req, err := rc.approvedStrategy.SetApproved(
				ApprovedStateReq{
					State: PaidStateTo(state),
					Token: TokenFromContext(ctx),
					UserID: lo.Ternary(
						user.IsAdminOrSuperAdmin(),
						mo.None[string](),
						mo.Some(user.ID),
					),
				},
			)
			if err != nil {
				return nil, rc.FormatError(
					codes.Internal,
					fmt.Errorf("can't set approved: %w", err),
				)
			}

			r.Query.Filter.And = append(r.Query.Filter.And, req)
		}
	}

	reports, err := rc.reportService.GetReports(
		ctx,
		r,
	)
	if err != nil {
		return nil, rc.HandlerError(err)
	}

	count, err := rc.reportService.CountReports(
		ctx,
		r,
	)
	if err != nil {
		return nil, rc.HandlerError(err)
	}

	return rc.marshaller.MarshalGetReportsPaginatedResp(
		req,
		reports,
		count,
	), nil
}

// Set report as paid
// Is report not found return REPORT_NOT_FOUND error
// Is report already paid return REPORT_ALREADY_PAID error
func (rc *ReportsController) PaidForReport(
	ctx context.Context,
	req *reportsgrpc.PaidForReportReq,
) (*reportsgrpc.PaidForReportResp, error) {
	user, err := rc.auth.AuthUser(
		ctx,
		TokenFromContext(ctx),
	)
	if err != nil {
		return nil, rc.HandlerError(err)
	}

	err = rc.reportService.SetReportPaid(
		ctx,
		reports.SetReportPaidReq{
			ID: req.ReportId,
			By: user,
		},
	)
	if err == reports.ErrReportNotFound {
		return &reportsgrpc.PaidForReportResp{
			Result: &reportsgrpc.PaidForReportResp_Error{
				Error: reportsgrpc.ReportsServiceErrors_REPORT_NOT_FOUND,
			},
		}, nil
	} else if err == reports.ErrReportsIsAlreadyPaid {
		return &reportsgrpc.PaidForReportResp{
			Result: &reportsgrpc.PaidForReportResp_Error{
				Error: reportsgrpc.ReportsServiceErrors_REPORT_ALREADY_PAID,
			},
		}, nil
	} else if errors.Is(err, reports.ErrCantSetReportPaid) {
		return nil, status.New(codes.PermissionDenied, err.Error()).Err()
	} else if err != nil {
		return nil, rc.HandlerError(err)
	}

	return &reportsgrpc.PaidForReportResp{
		Result: &reportsgrpc.PaidForReportResp_Empty{},
	}, nil
}
