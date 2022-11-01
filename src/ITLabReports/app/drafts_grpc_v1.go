package app

import (
	"github.com/RTUITLab/ITLab-Reports/transport/draft/grpc/servers/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
	pb "github.com/RTUITLab/ITLab/proto/reports/v1"
)

func (a *App) BuildDraftGRPCV1(e report.Endpoints) pb.DraftsServer {
	return servers.NewServer(
		e,
		servers.WithAuther(a.Auther),
	)
}
