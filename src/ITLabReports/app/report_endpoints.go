package app

import (
	"github.com/RTUITLab/ITLab-Reports/service/reports"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
)

func (a *App) BuildReportsEndpoints(s reports.Service) report.Endpoints {
	return report.MakeEndpoints(s)
}

