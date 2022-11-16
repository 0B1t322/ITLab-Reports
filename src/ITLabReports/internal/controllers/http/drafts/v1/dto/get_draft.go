package dto

type GetDraftReq struct {
	ID string `json:"id" uri:"id" validate:"required"`
}
