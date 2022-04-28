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