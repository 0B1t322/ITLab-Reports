package app

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/transport/report"
	"github.com/RTUITLab/ITLab-Reports/transport/report/http/servers/v1"
)

func (a *App) BuildReportsHTTPV1(
	e report.Endpoints,
	draftService servers.DraftService,
) {
	servers.NewServer(
		context.Background(),
		a.Router,
		e,
		servers.WithAuther(a.Auther),
		servers.WithDraftService(draftService),
	)
}