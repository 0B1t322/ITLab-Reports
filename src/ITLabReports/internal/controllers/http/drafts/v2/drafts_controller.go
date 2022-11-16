package drafts

import (
	"net/http"

	drafts "github.com/RTUITLab/ITLab-Reports/internal/domain/drafts/service"
	user "github.com/RTUITLab/ITLab-Reports/internal/domain/user/service"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

type DraftsController struct {
	draftService drafts.DraftsService
	auth         user.UserService

	marshaller DraftsMarshaller

	AuthErrorHandler
	ErrorFormatter
}

func NewDraftsController(
	draftService drafts.DraftsService,
	auth user.UserService,
) *DraftsController {
	formatter := ErrorFormatter{}
	return &DraftsController{
		draftService:     draftService,
		auth:             auth,
		ErrorFormatter:   formatter,
		AuthErrorHandler: AuthErrorHandler{Formatter: formatter},
		marshaller:       DraftsMarshaller{},
	}

}

func NewDraftsControllerFrom(i *do.Injector) (*DraftsController, error) {
	draftService := do.MustInvoke[drafts.DraftsService](i)
	auth := do.MustInvoke[user.UserService](i)

	return NewDraftsController(draftService, auth), nil
}

func (c *DraftsController) Build(r gin.IRouter) {
	draft := r.Group("/v2/draft")
	{
		draft.GET("", c.GetDraft)
	}
}

func (dc *DraftsController) HandlerError(c *gin.Context, err error) {
	if dc.AuthErrorHandler.HandlerError(c, err) {
		return
	}

	switch {
	default:
		dc.ErrorFormatter.FormatError(c, err, http.StatusInternalServerError)
		return
	}
}

// GetDrafts
//
// @Tags draft
//
// @Summary return drafts
//
// @Description.markdown get_drafts_v2
//
// @Router /reports/v2/draft [get]
//
// @Security ApiKeyAuth
//
// @Param offset query number false "offset"
//
// @Param limit query number false "limit"
//
// @Produce json
//
// @Success 200 {object} dto.GetDraftsResp
func (dc *DraftsController) GetDraft(c *gin.Context) {
	user, err := dc.auth.AuthUser(
		c,
		c.GetHeader("Authorization"),
	)
	if err != nil {
		dc.HandlerError(c, err)
		return
	}

	var req GetDraftsReq
	{
		if err := c.ShouldBindQuery(&req); err != nil {
			dc.FormatError(c, err, http.StatusBadRequest)
			return
		}
	}

	r := dc.marshaller.MarshallGetDraftsReq(req, user)

	drafts, err := dc.draftService.GetDrafts(c, r)
	if err != nil {
		dc.HandlerError(c, err)
		return
	}

	count, err := dc.draftService.CountDrafts(c, r)
	if err != nil {
		dc.HandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, dc.marshaller.MarshallGetDraftsResp(req, drafts, count))
}
