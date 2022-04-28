package dto

import (
	"context"
	"net/http"
)

type DeleteDraftResp struct {

}

func EncodeDeleteDraftResp(
	ctx context.Context,
	w http.ResponseWriter,
	resp * DeleteDraftResp,
) error {
	w.WriteHeader(http.StatusNoContent)
	w.Header().Add("Content-Type", "text/plain")
	return nil
}