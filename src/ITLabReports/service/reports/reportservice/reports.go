package reportservice

import (
	"context"
	"fmt"

	"github.com/RTUITLab/ITLab-Reports/aggragate/report"
	reportdomain "github.com/RTUITLab/ITLab-Reports/domain/report"
	"github.com/RTUITLab/ITLab-Reports/domain/report/mongo"
	"github.com/RTUITLab/ITLab-Reports/pkg/errors"
	service "github.com/RTUITLab/ITLab-Reports/service/reports"
	"github.com/samber/lo"
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

func WithMongoRepositoryAndCollectionName(
	connString string,
	collectionName string,
) ServiceConfiguration {
	return func(rp *ReportService) error {
		repo, err := mongo.New(
			context.Background(),
			connString,
			mongo.WithCollectionName(collectionName),
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

func (u UpdateReportToIReport) GetName() string {
	return u.Name.OrEmpty()
}

func (u UpdateReportToIReport) GetText() string {
	return u.Text.OrEmpty()
}

func (u UpdateReportToIReport) GetReporter() string {
	return ""
}

func (u UpdateReportToIReport) GetImplementer() string {
	return u.Implementer.OrEmpty()
}

func (u UpdateReportToIReport) Validate() error {
	var validatorOpts []report.ReportValidateOptions
	if u.Name.IsPresent() {
		validatorOpts = append(validatorOpts, report.WithValidateName())
	}

	if u.Text.IsPresent() {
		validatorOpts = append(validatorOpts, report.WithValidateText())
	}

	if u.Implementer.IsPresent() {
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
// 	catchable errors:
// 		ErrGetReportsBadParams as target
// Target errors catch by:
// 		errors.Is(err, ErrGetReportsBadParams)
func (r *ReportService) GetReports(
	ctx context.Context,
	params *reportdomain.GetReportsParams,
) ([]*report.Report, error) {
	if err := r.ValidateGetReportsParams(params); err != nil {
		return nil, errors.Wrap(err, service.ErrGetReportsBadParams)
	}

	return r.ReportRepository.GetReports(
		ctx,
		params,
	)
}

func (r *ReportService) ValidateGetReportsParams(
	params *reportdomain.GetReportsParams,
) error {
	// Validate sort params
	if sort := params.Filter.SortParams; len(sort) > 0 {
		// Check that in each item only one field is set
		// Check that not duplicated fields
		// Check that not empty sort item
		var (
			isDuplicatedFields    bool
			isMoreThanOneFieldSet bool
			isEmptyItem           bool
		)
		{
			duplicates := lo.FindDuplicates(
				lo.Map(
					sort,
					func(p reportdomain.GetReportsSort, _ int) string {
						if p.NameSort.IsPresent() && p.DateSort.IsPresent() {
							isMoreThanOneFieldSet = true
						} else if p.NameSort.IsAbsent() && p.DateSort.IsAbsent() {
							isEmptyItem = true
						}

						if p.NameSort.IsPresent() {
							return "name"
						}
						if p.DateSort.IsPresent() {
							return "date"
						}
						return ""
					},
				),
			)
			if len(duplicates) > 0 {
				isDuplicatedFields = true
			}
		}
		if isEmptyItem {
			return fmt.Errorf("empty sort argument")
		}
		if isMoreThanOneFieldSet {
			return fmt.Errorf("only one sort field can be set in slice item")
		}
		if isDuplicatedFields {
			return fmt.Errorf("you can't sort by field twice")
		}
	}
	return nil
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
	ctx context.Context,
	filter *reportdomain.GetReportsFilterFieldsWithOrAnd,
) (int64, error) {
	return r.ReportRepository.CountByFilter(
		ctx,
		filter,
	)
}
