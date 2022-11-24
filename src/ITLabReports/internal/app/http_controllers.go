package app

import (
	draftsv1 "github.com/RTUITLab/ITLab-Reports/internal/controllers/http/drafts/v1"
	draftsv2 "github.com/RTUITLab/ITLab-Reports/internal/controllers/http/drafts/v2"
	reportsv1 "github.com/RTUITLab/ITLab-Reports/internal/controllers/http/reports/v1"
	reportsv2 "github.com/RTUITLab/ITLab-Reports/internal/controllers/http/reports/v2"
	"github.com/samber/lo"

	"github.com/gin-gonic/gin"
)

type IHTTPController interface {
	Build(r gin.IRouter)
}

func (a *App) configureHTTPControllers() []IHTTPController {
	return []IHTTPController{
		lo.Must(reportsv1.NewReportsControllerFrom(a.injector)),
		lo.Must(reportsv2.NewReportsControllerFrom(a.injector)),
		lo.Must(draftsv1.NewDraftsControllerFrom(a.injector)),
		lo.Must(draftsv2.NewDraftsControllerFrom(a.injector)),
	}
}

func (a *App) ConfigureHTTPControllerOptions() {

}
