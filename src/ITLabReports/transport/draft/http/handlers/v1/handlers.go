package handlers

import (
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/transport/draft/http/dto/v1"
	endpoints "github.com/RTUITLab/ITLab-Reports/transport/draft/http/endpoints/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/http/options/serverbefore"
	errenc "github.com/RTUITLab/ITLab-Reports/transport/draft/http/options/servererrorencoder"

	genhttp "github.com/RTUITLab/ITLab-Reports/pkg/transport/http"

	kithttp "github.com/go-kit/kit/transport/http"
)

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
// @Success 200 {object} dto.GetDraftResp
func GetDraft(
	e endpoints.Endpoints,
) http.Handler {
	return genhttp.NewServer(
		e.GetDraft,
		dto.DecodeGetDraftReq,
		dto.EncodeGetDraftResp,
		kithttp.ServerBefore(
			serverbefore.TokenFromReq,
		),
		kithttp.ServerErrorEncoder(
			errenc.EncodeError,
		),
	)
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
// @Success 200 {object} dto.GetDraftsResp
func GetDrafts(
	e endpoints.Endpoints,
) http.Handler {
	return genhttp.NewServer(
		e.GetDrafts,
		dto.DecodeGetDraftsReq,
		dto.EncodeGetDraftsResp,
		kithttp.ServerBefore(
			serverbefore.TokenFromReq,
		),
		kithttp.ServerErrorEncoder(
			errenc.EncodeError,
		),
	)
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
func DeleteDraft(
	e endpoints.Endpoints,
) http.Handler {
	return genhttp.NewServer(
		e.DeleteDraft,
		dto.DecodeDeleteDraftReq,
		dto.EncodeDeleteDraftResp,
		kithttp.ServerBefore(
			serverbefore.TokenFromReq,
		),
		kithttp.ServerErrorEncoder(
			errenc.EncodeError,
		),
	)
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
// @Param body body dto.UpdateDraftReq true "a body"
//
// @Security ApiKeyAuth
//
// @Accept json
// 
// @Produce json
//
// @Success 200 {object} dto.UpdateDraftResp
func UpdateDraft(
	e endpoints.Endpoints,
) http.Handler {
	return genhttp.NewServer(
		e.UpdateDraft,
		dto.DecodeUpdateDraftReq,
		dto.EncodeUpdateDraftResp,
		kithttp.ServerBefore(
			serverbefore.TokenFromReq,
		),
		kithttp.ServerErrorEncoder(
			errenc.EncodeError,
		),
	)
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
// @Param body body dto.CreateDraftReq true "a body"
// 
// @Param implementer query string false "a id of implementer"
//
// @Security ApiKeyAuth
//
// @Accept json
// 
// @Produce json
//
// @Success 201 {object} dto.CreateDraftResp
func CreateDraft(
	e endpoints.Endpoints,
) http.Handler {
	return genhttp.NewServer(
		e.CreateDraft,
		dto.DecodeCreateDraftReq,
		dto.EncodeCreateDraftResp,
		kithttp.ServerBefore(
			serverbefore.TokenFromReq,
		),
		kithttp.ServerErrorEncoder(
			errenc.EncodeError,
		),
	)
}