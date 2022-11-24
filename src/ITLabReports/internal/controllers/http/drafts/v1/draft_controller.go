package drafts

import (
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/internal/common/errors"

	drafts "github.com/RTUITLab/ITLab-Reports/internal/domain/drafts/service"
	user "github.com/RTUITLab/ITLab-Reports/internal/domain/user/service"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
)

type DraftsController struct {
	draftsService drafts.DraftsService
	auth          user.UserService

	marshaller DraftsMarshaller

	ErrorFormatter
	AuthErrorHandler
}

func NewDraftsController(
	draftsService drafts.DraftsService,
	auth user.UserService,
) *DraftsController {
	formatter := ErrorFormatter{}
	c := &DraftsController{
		draftsService:    draftsService,
		auth:             auth,
		ErrorFormatter:   formatter,
		AuthErrorHandler: AuthErrorHandler{Formatter: formatter},
		marshaller:       DraftsMarshaller{},
	}

	return c
}

func NewDraftsControllerFrom(i *do.Injector) (*DraftsController, error) {
	draftsService := do.MustInvoke[drafts.DraftsService](i)
	auth := do.MustInvoke[user.UserService](i)

	return NewDraftsController(draftsService, auth), nil
}

func (c *DraftsController) Build(r gin.IRouter) {
	drafts := r.Group("/v1/draft")
	{
		drafts.GET("", c.GetDrafts)
		drafts.POST("", c.CreateDraft)
		drafts.GET("/:id", c.GetDraft)
		drafts.PUT("/:id", c.UpdateDraft)
		drafts.DELETE("/:id", c.DeleteDraft)
	}
}

func (dc *DraftsController) HandlerError(c *gin.Context, err error) {
	if dc.AuthErrorHandler.HandlerError(c, err) {
		return
	}

	switch {
	case errors.Is(err, aggregate.ErrDraftValidation):
		dc.FormatError(c, errors.Unwrap(err), http.StatusBadRequest)
		return
	case err == drafts.ErrDraftNotFound:
		dc.FormatError(c, err, http.StatusNotFound)
		return
	case errors.Is(err, drafts.ErrCantGetDraft),
		errors.Is(err, drafts.ErrCantDeleteDraft),
		errors.Is(err, drafts.ErrCantUpdateDraft):
		dc.FormatError(c, err, http.StatusForbidden)
		return
	default:
		logrus.WithFields(
			logrus.Fields{
				"controller": "drafts",
				"transport":  "http",
				"handler":    c.HandlerName(),
			},
		).Error(err)
		dc.FormatError(c, err, http.StatusInternalServerError)
		return
	}
}

// GetDraft
//
// @Tags draft
//
// @Summary return a draft
//
// @Description.markdown get_draft
//
// @Router /reports/v1/draft/{id} [get]
//
// @Security ApiKeyAuth
//
// @Param id path string true "id of draft"
//
// @Produce json
//
// @Success 200 {object} DraftView
func (dc *DraftsController) GetDraft(c *gin.Context) {
	user, err := dc.auth.AuthUser(
		c,
		c.GetHeader("Authorization"),
	)
	if err != nil {
		dc.HandlerError(c, err)
		return
	}

	var req GetDraftReq
	{
		if err := c.ShouldBindUri(&req); err != nil {
			dc.FormatError(c, err, http.StatusBadRequest)
			return
		}
	}

	draft, err := dc.draftsService.GetDraft(
		c,
		dc.marshaller.MarshallGetDraftReq(req, user),
	)
	if err != nil {
		dc.HandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, DraftViewFrom(draft))
}

// GetDrafts
//
// @Tags draft
//
// @Summary return a drafts for user
//
// @Description.markdown get_drafts
//
// @Router /reports/v1/draft [get]
//
// @Security ApiKeyAuth
//
// @Produce json
//
// @Success 200 {object} DraftsView
func (dc *DraftsController) GetDrafts(c *gin.Context) {
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
		//MARK: Add bind if needed
	}

	drafts, err := dc.draftsService.GetDrafts(
		c,
		dc.marshaller.MarshallGetDraftsReq(req, user),
	)
	if err != nil {
		dc.HandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, DraftsViewFrom(drafts))
}

// DeleteDraft
//
// @Tags draft
//
// @Summary delete a draft
//
// @Description.markdown delete_draft
//
// @Router /reports/v1/draft/{id} [delete]
//
// @Param id path string true "id of draft"
//
// @Security ApiKeyAuth
//
// @Produce json,plain
//
// @Success 204
func (dc *DraftsController) DeleteDraft(c *gin.Context) {
	user, err := dc.auth.AuthUser(
		c,
		c.GetHeader("Authorization"),
	)
	if err != nil {
		dc.HandlerError(c, err)
		return
	}

	var req DeleteDraftReq
	{
		if err := c.ShouldBindUri(&req); err != nil {
			dc.FormatError(c, err, http.StatusBadRequest)
			return
		}
	}

	err = dc.draftsService.DeleteDraft(
		c,
		dc.marshaller.MarshallDeleteDraftReq(req, user),
	)
	if err != nil {
		dc.HandlerError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateDraft
//
// @Tags draft
//
// @Summary update a draft
//
// @Description.markdown update_draft
//
// @Router /reports/v1/draft/{id} [put]
//
// @Param id path string true "id of draft"
//
// @Param body body UpdateDraftReq true "a body"
//
// @Security ApiKeyAuth
//
// @Accept json
//
// @Produce json
//
// @Success 200 {object} DraftView
func (dc *DraftsController) UpdateDraft(c *gin.Context) {
	user, err := dc.auth.AuthUser(
		c,
		c.GetHeader("Authorization"),
	)
	if err != nil {
		dc.HandlerError(c, err)
		return
	}

	var req UpdateDraftReq
	{
		if err := c.ShouldBindUri(&req); err != nil {
			dc.FormatError(c, err, http.StatusBadRequest)
			return
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			dc.FormatError(c, err, http.StatusBadRequest)
			return
		}
	}

	draft, err := dc.draftsService.UpdateDraft(
		c,
		dc.marshaller.MarshallUpdateDraftReq(req, user),
	)
	if err != nil {
		dc.HandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, DraftViewFrom(draft))
}

// CreateDraft
//
// @Tags draft
//
// @Summary create a draft
//
// @Description.markdown create_draft
//
// @Router /reports/v1/draft [post]
//
// @Param body body CreateDraftReq true "body"
//
// @Param implementer query string false "a id of implementer"
//
// @Security ApiKeyAuth
//
// @Accept json
//
// @Produce json
//
// @Success 201 {object} DraftView
func (dc *DraftsController) CreateDraft(c *gin.Context) {
	user, err := dc.auth.AuthUser(
		c,
		c.GetHeader("Authorization"),
	)
	if err != nil {
		dc.HandlerError(c, err)
		return
	}

	var req CreateDraftReq
	{
		if err := c.ShouldBindJSON(&req); err != nil {
			dc.FormatError(c, err, http.StatusBadRequest)
			return
		}

		if err := c.ShouldBindQuery(&req); err != nil {
			dc.FormatError(c, err, http.StatusBadRequest)
			return
		}
	}

	draft, err := dc.draftsService.CreateDraft(
		c,
		dc.marshaller.MarshallCreateDraftReq(req, user),
	)
	if err != nil {
		dc.HandlerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, DraftViewFrom(draft))
}
