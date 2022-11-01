package handlers

import (
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/transport/draft/http/dto/v2"
	"github.com/RTUITLab/ITLab-Reports/transport/draft/http/endpoints/v2"
	"github.com/RTUITLab/ITLab-Reports/transport/http/options/serverbefore"

	genhttp "github.com/RTUITLab/ITLab-Reports/pkg/transport/http"

	kithttp "github.com/go-kit/kit/transport/http"

	errenc "github.com/RTUITLab/ITLab-Reports/transport/draft/http/options/servererrorencoder"
)

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
func GetDrafts(
	e endpoints.Endpoints,
) http.Handler {
	return genhttp.NewServer(
		e.GetDrafts,
		dto.DecodeGetDraftsReq,
		dto.EncodeGetDraftResp,
		kithttp.ServerBefore(
			serverbefore.TokenFromReq,
		),
		kithttp.ServerErrorEncoder(
			errenc.EncodeError,
		),
	)
}
