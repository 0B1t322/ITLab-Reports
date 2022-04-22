package reportservice

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/aggragate/report"
	reportdomain "github.com/RTUITLab/ITLab-Reports/domain/report"
	"github.com/RTUITLab/ITLab-Reports/domain/report/mongo"
	"github.com/RTUITLab/ITLab-Reports/pkg/errors"
	service "github.com/RTUITLab/ITLab-Reports/service/reports"
	m "go.mongodb.org/mongo-driver/mongo"
)

type ServiceConfiguration func(rp *ReportService) error

type ReportService struct {
	ReportRepository reportdomain.ReportRepository
}

func New(
	cfgs ...ServiceConfiguration,
) (*ReportService, error) {
	service := &ReportService{}

	for _, cfg := range cfgs {
		if err := cfg(service); err != nil {
			return nil, err
		}
	}

	return service, nil
}

func WithMongoRepository(connString string) ServiceConfiguration {
	return func(rp *ReportService) error {
		repo, err := mongo.New(
			context.Background(),
			connString,
		)
		if err != nil {
			return err
		}

		rp.ReportRepository = repo
		return nil
	}
}

func WithMongoRepositoryWithClient(
	connString string,
	client *m.Client,
) ServiceConfiguration {
	return func(rp *ReportService) error {
		repo, err := mongo.New(
			context.Background(),
			connString,
			mongo.WithClient(client),
		)
		if err != nil {
			return err
		}

		rp.ReportRepository = repo
		return nil
	}
}

// GetReport return report by id
// 	catchable errors:
// 		ErrReportIDNotValid
// 		ErrReportNotFound
func (r *ReportService) GetReport(ctx context.Context, id string) (*report.Report, error) {
	report, err := r.ReportRepository.GetReport(ctx, id)
	switch {
	case err == reportdomain.ErrIDIsNotValid:
		return nil, service.ErrReportIDNotValid
	case err == reportdomain.ErrReportNotFound:
		return nil, service.ErrReportNotFound
	case err != nil:
		return nil, err
	}

	return report, nil
}

// DeleteReport delete report by id
// 	catchable errors:
// 		ErrReportIDNotValid
// 		ErrReportNotFound
func (r *ReportService) DeleteReport(ctx context.Context, id string) error {
	err := r.ReportRepository.DeleteReport(
		ctx,
		id,
	)
	switch {
	case err == reportdomain.ErrIDIsNotValid:
		return service.ErrReportIDNotValid
	case err == reportdomain.ErrReportNotFound:
		return service.ErrReportNotFound
	case err != nil:
		return err
	}

	return nil
}

type UpdateReportToIReport reportdomain.UpdateReportParams

func(u UpdateReportToIReport) GetName() string {
	return u.Name.MustGetValue()
}

func (u UpdateReportToIReport) GetText() string {
	return u.Text.MustGetValue()
}

func (u UpdateReportToIReport) GetReporter() string {
	return ""
}

func (u UpdateReportToIReport) GetImplementer() string {
	return u.Implementer.MustGetValue()
}

func (u UpdateReportToIReport) Validate() error {
	var validatorOpts []report.ReportValidateOptions
	if u.Name.HasValue() {
		validatorOpts = append(validatorOpts, report.WithValidateName())
	}

	if u.Text.HasValue() {
		validatorOpts = append(validatorOpts, report.WithValidateText())
	}

	if u.Implementer.HasValue() {
		validatorOpts = append(validatorOpts, report.WithValidateImplementor())
	}

	if len(validatorOpts) != 0 {
		return report.NewReportValidator(validatorOpts...).Validate(u)
	}

	return nil
}

// UpdateReport update reports by id and not nil optionals
// Name, Text, Implemtner fields can't be empty
// 	catchable errors:
// 		ErrReportIDNotValid
// 		ErrReportNotFound
// 		ErrValidationError as target
// Target errors catch by:
// 		errors.Is(err, ErrValidationError)
func (r *ReportService) UpdateReport(
	ctx context.Context, 
	id string, 
	params reportdomain.UpdateReportParams,
) (*report.Report, error) {
	if err := report.NewReportValidator(
		report.WithSelfValidator(),
	).Validate(UpdateReportToIReport(params)); err != nil {
		return nil, errors.Wrap(err, service.ErrValidationError)
	}

	report, err := r.ReportRepository.UpdateReport(
		ctx,
		id,
		params,
	)
	switch {
	case err == reportdomain.ErrIDIsNotValid:
		return nil, service.ErrReportIDNotValid
	case err == reportdomain.ErrReportNotFound:
		return nil, service.ErrReportNotFound
	case err != nil:
		return nil, err
	}

	return report, nil
}

// GetReports return reports acording to filters
// 	don't have catchable errors
func (r *ReportService) GetReports(
	ctx context.Context, 
	params *reportdomain.GetReportsParams,
) ([]*report.Report, error) {
	return r.ReportRepository.GetReports(
		ctx,
		params,
	)
}

// CreateReport create report and return it
// 	catchable errors:
// 		ErrValidationError as target
// Target errors catch by:
// 		errors.Is(err, ErrValidationError)
func (r *ReportService) CreateReport(
	ctx context.Context, 
	model *report.Report,
) (*report.Report, error) {
	if err := report.NewReportValidator().Validate(model); err != nil {
		return nil, errors.Wrap(err, service.ErrValidationError)
	}

	created, err := r.ReportRepository.CreateReport(
		ctx,
		model,
	)
	if err != nil {
		return nil, err
	}

	return created, nil
}


// CountReport count report according to filter and return count
// don't have catchable errors
func (r *ReportService) CountReports(
	ctx	context.Context,
	filter *reportdomain.GetReportsFilterFieldsWithOrAnd,
) (int64, error) {
	return r.ReportRepository.CountByFilter(
		ctx,
		filter,
	)
}


