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
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNoContent)
	return nil
}