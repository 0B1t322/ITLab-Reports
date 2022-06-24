package app

import (
	"github.com/RTUITLab/ITLab-Reports/service/reports/reportservice"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
)

func (a *App) BuildReportService() error {
	service, err := reportservice.New(
		reportservice.WithMongoRepositoryAndCollectionName(
			a.cfg.MongoDB.URI,
			"reports",
		),
	)
	if err != nil {
		return err
	}

	a.ReportService = service

	return nil
}

func (a *App) BuildReportEndpoints() {
	a.ReportEndpoints = report.MakeEndpoints(a.ReportService)
}