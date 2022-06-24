package app

import (
	"github.com/RTUITLab/ITLab-Reports/service/reports/reportservice"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
)

func (a *App) BuildDraftService() error {
	service, err := reportservice.New(
		reportservice.WithMongoRepositoryAndCollectionName(
			a.cfg.MongoDB.URI,
			"drafts",
		),
	)
	if err != nil {
		return err
	}

	a.DraftService = service

	return nil
}

func (a *App) BuildDraftEndpoints() {
	a.DraftEndpoints = report.MakeEndpoints(a.DraftService)
}