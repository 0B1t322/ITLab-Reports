package app

import (
	"github.com/RTUITLab/ITLab-Reports/service/reports"
	"github.com/RTUITLab/ITLab-Reports/service/reports/reportservice"
)

func (a *App) BuildReportService() (reports.Service ,error) {
	service, err := reportservice.New(
		reportservice.WithMongoRepository(
			a.cfg.MongoDB.URI,
		),
	)
	if err != nil {
		return nil, err
	}

	return service, nil
}