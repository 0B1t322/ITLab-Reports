package app

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/transport/draft/http/endpoints/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/draft/http/servers/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
)

type DraftEndpoints = endpoints.Endpoints

func (a *App) BuildDraftsHTTPV1(
	e report.Endpoints,
) DraftEndpoints {
	return servers.NewServer(
		context.Background(),
		a.Router,
		e,
		servers.WithAuther(a.Auther),
		servers.WithIdChecker(a.IdChecker),
	)
}