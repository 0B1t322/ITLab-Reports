package app

import (
	"github.com/RTUITLab/ITLab-Reports/transport/report"
	"github.com/RTUITLab/ITLab-Reports/transport/report/grpc/servers/v1"
	pb "github.com/RTUITLab/ITLab/proto/reports/v1"
)

func (a *App) BuildReportsGRPCV1(
	e report.Endpoints,
) pb.ReportsServer {
	return servers.NewServer(
		e,
		servers.WithAuther(a.Auther),
	)
}