package app

import (
	draftsv1 "github.com/RTUITLab/ITLab-Reports/internal/controllers/grpc/drafts/v1"
	reportsv1 "github.com/RTUITLab/ITLab-Reports/internal/controllers/grpc/reports/v1"
	"github.com/samber/lo"
	"google.golang.org/grpc"
)

type GRPCController interface {
	Build(r *grpc.Server)
}

func (a *App) configureGRPCControllers() []GRPCController {
	return []GRPCController{
		lo.Must(reportsv1.NewReportsControllerFrom(a.injector)),
		lo.Must(draftsv1.NewDraftControllerFrom(a.injector)),
	}
}

func (a *App) ConfigureGRPCControllerOptions() {
}
