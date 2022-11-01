package app

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/transport/draft/http/servers/v2"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
)

func (a *App) BuildDraftsHTTPV2(
	e report.Endpoints,
) {
	servers.NewServer(
		context.Background(),
		a.Router,
		e,
		servers.WithAuther(a.Auther),
	)
}
