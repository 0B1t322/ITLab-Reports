package app

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/pkg/adapters/toapprovereportsidgetter"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
	"github.com/RTUITLab/ITLab-Reports/transport/report/http/servers/v2"
)

func (a *App) BuildReportsHTTPV2(e report.Endpoints) {
	servers.NewServer(
		context.Background(),
		a.Router,
		e,
		servers.WithAuther(a.Auther),
		servers.WithApprovedreportsIdGetter(toapprovereportsidgetter.ToApproveReportsIdGetter(a.SalaryService)),
	)
}
